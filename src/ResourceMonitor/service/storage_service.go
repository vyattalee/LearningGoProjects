package service

import (
	"fmt"
	"github.com/LearningGoProjects/ResourceMonitor/pb"
	"github.com/shirou/gopsutil/disk"
)

func getStorageInfo() (*pb.Storage, error) {

	partitionsInfos, err1 := disk.Partitions(true) //所有分区
	//fmt.Println(partitionsInfos)
	//usageInfo, err2 := disk.Usage("~/") //指定某路径的硬盘使用情况
	//fmt.Println(usageInfo)
	//ioCountinfo, err3 := disk.IOCounters() //所有硬盘的io信息
	//fmt.Println(ioCountinfo)
	var err2, err3 error //
	var usageInfo *disk.UsageStat
	//var ioCountinfo map[string]disk.IOCountersStat

	var partitions []*pb.PartitionStat
	partitions = make([]*pb.PartitionStat, len(partitionsInfos))
	var diskUsage []*pb.UsageStat
	diskUsage = make([]*pb.UsageStat, len(partitionsInfos))
	//var ioCount []map[string]disk.IOCountersStat
	//ioCount = make([]map[string]disk.IOCountersStat, len(partitionsInfos))

	for i, partitionsInfo := range partitionsInfos {
		partitions[i] = &pb.PartitionStat{
			Device:     partitionsInfo.Device,
			MountPoint: partitionsInfo.Mountpoint,
			FsType:     partitionsInfo.Fstype,
			Opts:       partitionsInfo.Opts,
		}

		usageInfo, err2 = disk.Usage(partitions[i].Device) //指定某路径的硬盘使用情况

		if usageInfo != nil {
			diskUsage[i] = &pb.UsageStat{
				Path:              usageInfo.Path,
				FsType:            usageInfo.Fstype,
				Total:             usageInfo.Total,
				Free:              usageInfo.Free,
				Used:              usageInfo.Used,
				UsedPercent:       usageInfo.UsedPercent,
				InodesTotal:       usageInfo.InodesTotal,
				InodesUsed:        usageInfo.InodesUsed,
				InodesFree:        usageInfo.InodesFree,
				InodesUsedPercent: usageInfo.InodesUsedPercent,
			}
		}
		//ioCountinfo, err3  = disk.IOCounters(partitions[i].Device)
		//
		//ioCount[i] = ioCountinfo

	}

	return &pb.Storage{Partition: partitions, Usage: diskUsage /*, IoCount:ioCount */}, fmt.Errorf("%v %v %v", err1, err2, err3)
	//[{"device":"C:","mountpoint":"C:","fstype":"NTFS","opts":"rw.compress"} {"device":"D:","mountpoint":"D:","fstype":"NTFS","opts":"rw.compress"} {"device":"E:","mountpoint":"E:","fstype":"NTFS","opts":"rw.compress"} ]
	//{"path":"E:","fstype":"","total":107380965376,"free":46790828032,"used":60590137344,"usedPercent":56.425398236866755,"inodesTotal":0,"inodesUsed":0,"inodesFree":0,"inodesUsedPercent":0}
	//map[C::{"readCount":0,"mergedReadCount":0,"writeCount":0,"mergedWriteCount":0,"readBytes":0,"writeBytes":4096,"readTime":0,"writeTime":0,"iopsInProgress":0,"ioTime":0,"weightedIO":0,"name":"C:","serialNumber":"","label":""} 。。。]

}
