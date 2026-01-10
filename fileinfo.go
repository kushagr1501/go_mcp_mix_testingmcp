package main

func GetFileInfo(client *MCPClient, path string) (*FILEINFO,error){
	data,err:=client.ToolCall("get_file_info", map[string]any{}{ 
		"path":path
	})

	if err != nil {
		return nil, err
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}

	info := &FileInfo{
		Path: path,
	}

	return info,nil




}