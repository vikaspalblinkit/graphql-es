// internal/elastic/client.go
package elastic

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"

	esv8 "github.com/elastic/go-elasticsearch/v8"
)

type Client struct {
	Client *esv8.Client
	Index  string
}

func (es *Client) CreateIndex(mappingPath string) error {
	body, err := os.ReadFile(mappingPath)
	if err != nil {
		return err
	}
	res, err := es.Client.Indices.Create(es.Index, es.Client.Indices.Create.WithBody(bytes.NewReader(body)))
	if err != nil || res.IsError() {
		return fmt.Errorf("create index error: %v", err)
	}
	return nil
}

func (es *Client) BulkUploadCSV(folder string) error {
	files, _ := filepath.Glob(filepath.Join(folder, "*.csv"))
	var bulkBuffer bytes.Buffer

	for _, file := range files {
		csvFile, err := os.Open(file)
		if err != nil {
			return err
		}
		reader := csv.NewReader(csvFile)
		_, _ = reader.Read() // read header, ignore variable
		for {
			record, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}
			// Map CSV columns to struct fields
			var salary float64
			if s := record[4]; s != "" {
				salary, _ = strconv.ParseFloat(s, 64)
			}
			data := Compensation{
				Timestamp:     record[0],
				AgeRange:      record[1],
				Industry:      record[2],
				JobTitle:      record[3],
				AnnualSalary:  salary,
				Currency:      record[5],
				Location:      record[6],
				Experience:    record[7],
				JobContext:    record[8],
				OtherCurrency: record[9],
			}
			meta := map[string]map[string]string{"index": {"_index": es.Index}}
			metaJSON, _ := json.Marshal(meta)
			dataJSON, _ := json.Marshal(data)
			bulkBuffer.Write(metaJSON)
			bulkBuffer.WriteString("\n")
			bulkBuffer.Write(dataJSON)
			bulkBuffer.WriteString("\n")
		}
		csvFile.Close()
	}
	res, err := es.Client.Bulk(bytes.NewReader(bulkBuffer.Bytes()), es.Client.Bulk.WithIndex(es.Index))
	if err != nil || res.IsError() {
		return fmt.Errorf("bulk upload error: %v", err)
	}
	return nil
}

func (es *Client) Search(filters map[string]interface{}, sortBy *string, sortOrder string, fields []string, limit int) ([]*Compensation, error) {
	query := map[string]interface{}{
		"query": map[string]interface{}{"bool": map[string]interface{}{"must": []interface{}{}}},
	}
	must := query["query"].(map[string]interface{})["bool"].(map[string]interface{})["must"].([]interface{})
	for k, v := range filters {
		if k == "salary_gte" {
			must = append(must, map[string]interface{}{"range": map[string]interface{}{"annual_salary": map[string]interface{}{"gte": v}}})
		} else if k == "id" { // changed from _id to id
			must = append(must, map[string]interface{}{"ids": map[string]interface{}{"values": []interface{}{v}}})
		} else {
			must = append(must, map[string]interface{}{"match": map[string]interface{}{k: v}})
		}
	}
	query["query"].(map[string]interface{})["bool"].(map[string]interface{})["must"] = must
	if sortBy != nil && *sortBy != "" {
		query["sort"] = []map[string]interface{}{{*sortBy: map[string]string{"order": sortOrder}}}
	}
	if len(fields) > 0 {
		query["_source"] = fields
	}
	if limit > 0 {
		query["size"] = limit
	}
	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(query)
	res, err := es.Client.Search(es.Client.Search.WithIndex(es.Index), es.Client.Search.WithBody(&buf))
	if err != nil || res.IsError() {
		return nil, fmt.Errorf("search error: %v", err)
	}
	defer res.Body.Close()
	var r map[string]interface{}
	json.NewDecoder(res.Body).Decode(&r)
	hits := r["hits"].(map[string]interface{})["hits"].([]interface{})
	var results []*Compensation
	for _, h := range hits {
		src := h.(map[string]interface{})["_source"]
		data, _ := json.Marshal(src)
		var c Compensation
		json.Unmarshal(data, &c)
		c.ID = h.(map[string]interface{})["_id"].(string)
		results = append(results, &c)
	}
	return results, nil
}

func (es *Client) GetByID(id string) (*Compensation, error) {
	res, err := es.Client.Get(es.Index, id)
	if err != nil || res.IsError() {
		return nil, fmt.Errorf("not found")
	}
	defer res.Body.Close()
	var r map[string]interface{}
	json.NewDecoder(res.Body).Decode(&r)
	var c Compensation
	data, _ := json.Marshal(r["_source"])
	json.Unmarshal(data, &c)
	c.ID = id // set the ID field from the argument
	return &c, nil
}
