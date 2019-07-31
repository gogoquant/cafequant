import { Component, OnInit } from "@angular/core";
import { ClusterService } from "./cluster.service";

@Component({
  selector: "app-cluster",
  templateUrl: "./cluster.component.html",
  styleUrls: ["./cluster.component.scss"]
})
export class ClusterComponent implements OnInit {
  clustersListData = [];
  pageIndex = 1;
  pageSize = 10;
  total = 1;
  loading = true;
  keyword = "";
  constructor(private clusterService: ClusterService) {}
  ngOnInit() {
    this.clusterData();
  }

  // 根据集群名称查询结果
  searchValue(): void {
    let new_list = [];
    if (this.keyword == "" || this.keyword.indexOf(" ") == 0) {
      this.clusterData();
    } else {
      for (const item of this.clustersListData) {
        if (item.name.indexOf(this.keyword.trim()) >= 0) {
          new_list.push(item);
        }
      }
      this.clustersListData = new_list;
    }
  }

  // 获取集群列表数据
  clusterData(): void {
    this.clusterService
      .getClusters(this.pageIndex, this.pageSize)
      .subscribe((data: any) => {
        this.loading = false;
        this.total = data.total;
        this.clustersListData = data.data;
      });
  }
}
