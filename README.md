# ArcherNetDisk
开源的，高可用，快速上传和下载，多链路工作，去重，支持存储和调度节点的平行扩展（依赖etcd）的网盘存储项目，主要面向日常用户使用
项目架构 v1:
                                                       
  client --------------req------------- master (standalone or cluster by nginx load balance)  
    | <------ --------resp------------- | /|\   --allocate space from worker load in local cache ----|
    |                                   |  |    --update local load cache by pulling from etcd   ----|
    |                                   |  |    --log complited jobs(include paused) state       ----| 
    |                                   |  |                                                         |---->etcd
    |                        distribute |  | sync                                                    |
    |                            tasks \|/ | state                                                   | 
    |---transmit-- file--splits------->workers   --------------add new worker nodes ----------- ---- | 
                       
                           
支持多种语言，提供一个Go语言编写的命令行工具

使用示例：







