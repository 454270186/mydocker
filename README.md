## MyDocker
![CNCF](https://img.shields.io/badge/CNCF-231F20.svg?style=for-the-badge&logo=CNCF&logoColor=white)
![docker](https://img.shields.io/badge/Docker-2496ED.svg?style=for-the-badge&logo=Docker&logoColor=white)

模仿Docker默认运行时runc实现的底层容器运行时

### Features
基于Linux Namespace, Cgroups, Overlayfs实现容器的基本操作：
    
- 容器资源限制(CPU核心数, CPU使用率, 内存使用率)
- 容器隔离
    - 容器进程隔离
    - 容器挂载点隔离
    - 容器文件读写隔离
    - 容器Rootfs隔离

- 容器外部卷轴挂载

- 容器后台运行(detach)

- 查看容器info和log

- 进入容器内部