package storage

//简明的负载信息和网络信息

type LoadInfo struct {
	Cpu
	Mem
	Disk
	Net
}

type Cpu struct {
	CpuNum        int     `json:"cpu_num"`        //cpu数量
	CpuPercentage float64 `json:"cpu_percentage"` //cpu当前使用量占比
	ThreadNum     int     `json:"thread_num"`     //当前线程数量
	CpuSpeed      int     `json:"cpu_speed"`      //cpu速度
}
type Mem struct {
	MemTotal int `json:"mem_total"` //内存总量
	MemUsed  int `json:"mem_used"`  //内存已使用量
	MemCache int `json:"mem_cache"` //内存缓冲大小
	MemSpeed int `json:"mem_speed"` //内存速度
}
type Disk struct {
	DiskTotal      int `json:"disk_total"`       //磁盘总量
	DiskUsed       int `json:"disk_used"`        //磁盘已使用量
	DiskReadSpeed  int `json:"disk_read_speed"`  //磁盘的读取速度
	DiskWriteSpeed int `json:"disk_write_speed"` //磁盘写入速度
	DiskRespTime   int `json:"disk_resp_time"`   //磁盘平均响应时间

}

//使用tcp传输
type Net struct {
	NetInterfaces []string `json:"net_interfaces"` //已经打开的网络 IP:端口数量 格式 IP:Port
}

func GetLoadInfo() *LoadInfo {
	//todo
	return nil
}
