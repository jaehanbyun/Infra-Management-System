# Infra Repo - CloudChain

<div align="center">
<img width="150" height="90" alt="스크린샷 2023-10-16 오전 1 38 58" src="https://github.com/Cloud-Chain/infra-repo/assets/80397512/a8c6dfd2-3752-4ad4-8c35-ff58b553434b">
</div>

## :thought_balloon: 개발 배경

부산대 2023 전기 졸업과제 진행 중 오픈스택에 구축된 쿠버네티스 환경에 애플리케이션 배포를 수행하였었습니다. 그러나, 오픈스택 상에 쿠버네티스 환경을 구축하는 절차가 비용이 많이 들어 이에 대한 부담을 자동화를 통해 줄이고자 인프라 관리 시스템을 기획하였습니다.

## :floppy_disk: 시스템 아키텍쳐

<img width="484" alt="image" src="https://github.com/Cloud-Chain/infra-repo/assets/80397512/914dc090-25bd-4fac-adec-443e2e65dfcd">

## :page_facing_up: 기능

#### 메인 페이지
<img width="455" alt="image" src="https://github.com/jaehanbyun/Infra-Management-System/assets/80397512/67e2139a-0c87-4c52-bf8a-0c1d7c5eb142">

- 오픈스택 플랫폼의 Instance, CPU, RAM, Storage에 대해 사용 가능한 양과 현재 사용 중인 양과 클러스터 배포 현황을 확인할 수 있습니다.

#### 클러스터 조회
<img width="455" alt="image" src="https://github.com/jaehanbyun/Infra-Management-System/assets/80397512/2a59bd81-b523-43a9-bef9-156797f9ac6d">

- 클러스터 조회 페이지에서는 구축된 클러스터의 이름, Node들의 Image, Node의 하드웨어(Flavor) 명세, Master Node와 Worker Node의 수와 현재 상태, 모니터링 대시보드의 접근 주소를 보여줍니다.

#### 클러스터 생성
<img width="413" alt="image" src="https://github.com/jaehanbyun/Infra-Management-System/assets/80397512/2be12a4f-e4a9-42a6-b1f5-af7ffde39eee">

- 사용자는 클러스터의 이름, Master Node와 Worker Node의 수, Node의 Image, Node의 하드웨어 명세(CPU, RAM, Disk)를 입력하고 CREATE 버튼을 누릅니다.
- 이후, 인프라 관리 시스템이 입력한 명세를 바탕으로 Jenkins 클러스터 배포 자동화 파이프라인을 통해 클러스터를 생성하게 된니다.

##### Jenkins 파이프라인 구동 화면 
<img width="455" alt="image" src="https://github.com/jaehanbyun/Infra-Management-System/assets/80397512/5f536e6b-14c9-41f7-880d-90d7abbaceb0">

Jenkins 파이프라인의 단계는 다음과 같습니다.
1. Setup Environment Variables
2. Terraform Init
3. Terraform Apply
4. Setup Kubespray On Bastion
5. Setup Kubectl on Bastion and configure Cluster Access
6. Install Node Exporter on All Nodes
7. Setup Prometheus & Grafana on Bastion Node

#### 클러스터 모니터링
<img width="455" alt="image" src="https://github.com/jaehanbyun/Infra-Management-System/assets/80397512/903355db-fe26-484b-b787-8b4304f2ca58">

- 모니터링 대시보드 접근 주소를 클릭하면, 해당하는 클러스터 모니터링 대시보드 화면이 나타나게 됩니다.
- 클러스터 모니터링 대시보드 화면에서는 각 Node 상태, CPU, Memory, Disk I/O, Network 트래픽 등의 클러스터 운영에 주요한 성능 지표들을 확인할 수 있습니다.

## :sparkles: Tech Stack
Kubespray, Terraform, Jenkins, Prometheus, Grafana, Go, React
