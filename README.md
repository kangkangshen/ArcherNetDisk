# ArcherNetDisk
开源的，高可用，快速上传和下载，多链路工作，去重，支持存储和调度节点的平行扩展（依赖etcd）的网盘存储项目，主要面向日常用户使用\n
项目架构 v1:\n
                                                       
  client --------------------req----------------------- master (standalone or cluster by nginx load balance)  \n
    | <----------------------resp----------------------- | /|\   --allocate space from worker load in local cache ----|\n
    |                                                    |  |    --update local load cache by pulling from etcd   ----|\n
    |                                                    |  |    --log complited jobs(include paused) state       ----| \n
    |                                                    |  |                                                         |---->etcd\n
    |                                         distribute |  | sync                                                    |\n
    |                                             tasks \|/ | state                                                   |\n 
    |------------transmit--- file----splits----------->workers   ----------------add new worker nodes --------------- |\n 
                       
                           
支持多种语言，提供一个Go语言编写的命令行工具

使用示例：







