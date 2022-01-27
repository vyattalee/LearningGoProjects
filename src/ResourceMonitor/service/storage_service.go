package service

import (
	"fmt"
	"github.com/LearningGoProjects/ResourceMonitor/pb"
	"github.com/shirou/gopsutil/disk"
)

func getStorageInfo() (*pb.Storage, error) {

	partitionsInfo, err1 := disk.Partitions(true) //所有分区
	fmt.Println(partitionsInfo)
	usageInfo, err2 := disk.Usage("E:") //指定某路径的硬盘使用情况
	fmt.Println(usageInfo)
	ioCountinfo, err3 := disk.IOCounters() //所有硬盘的io信息
	fmt.Println(ioCountinfo)

	return &pb.Storage{ /*Partition: partitionsInfo, Usage:usageInfo, IoCount:ioCountinfo*/ }, fmt.Errorf("%s %s %s", err1, err2, err3)
	//[{"device":"C:","mountpoint":"C:","fstype":"NTFS","opts":"rw.compress"} {"device":"D:","mountpoint":"D:","fstype":"NTFS","opts":"rw.compress"} {"device":"E:","mountpoint":"E:","fstype":"NTFS","opts":"rw.compress"} ]
	//{"path":"E:","fstype":"","total":107380965376,"free":46790828032,"used":60590137344,"usedPercent":56.425398236866755,"inodesTotal":0,"inodesUsed":0,"inodesFree":0,"inodesUsedPercent":0}
	//map[C::{"readCount":0,"mergedReadCount":0,"writeCount":0,"mergedWriteCount":0,"readBytes":0,"writeBytes":4096,"readTime":0,"writeTime":0,"iopsInProgress":0,"ioTime":0,"weightedIO":0,"name":"C:","serialNumber":"","label":""} 。。。]

}
