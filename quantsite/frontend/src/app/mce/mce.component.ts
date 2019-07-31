import { Component, OnInit } from '@angular/core';
import { LoggerService } from '../commons/services/logger.service';
import { MceService } from '../mce/mce.service';
import {PodInfo} from '../mce/mce.type';

@Component({
    selector: 'app-mce',
    templateUrl: '../mce/mce.component.html',
    styleUrls: ['../mce/mce.component.scss']
})


export class MceComponent implements OnInit {

    //  详情使用
    detailName = '';
    detailData = null;
    showDetail = false;

    //  分页使用
    pageNum = 1;
    pageSize = 10;
    total = 1;

    // 加载标志
    namespaceloading  = false;
    clustersloading = false;
    podsloading = false;
    apploading = false;

    // 选定记录
    namespace = '';
    app = '';
    cluster = '';


    namespaceData = [];
    clustersData = [];
    appsData = [];
    podsData = [];

    constructor(private logger: LoggerService,
              private mceService: MceService) { }

    ngOnInit() {
      console.log('初始化页面');
      this.searchData(false);
    }

    detailClick = (value) => {
        this.detailName = value;
        this.showDetail = true;
        for (const item of this.podsData) {
            if (item.name === value) {
                console.log('match success');
                this.detailData = item;
                break;
            }
        }

    }

    listClick = () => {
        this.showDetail = false;
        this.detailName = '';
        this.detailData = null;
    }

  // cluster change
  clusterChange( value: string): void {
    console.log(value);
    this.app = '';
    this.namespace = '';
    this.cluster = value;
    this.searchData();
  }

  // namespace change
  namespaceChange( value: string): void {
    console.log(value);
    this.app = '';
    this.namespace = value;
    this.searchData();
  }

  // namespace change
  appChange( value: string): void {
    console.log(value);
    this.app = value;
    this.searchData();
  }

  /* 获取 app 列表*/
  queryApps(): void {
    console.log('尝试获取app');
    if (this.cluster === '') {
      return;
    }

    if (this.namespace === '') {
      return;
    }

    if (this.apploading === true) {
      return;
    }
    this.apploading = true;
    this.mceService.getApps(this.cluster, this.namespace)
      .subscribe((data: any) => {
        if (data !== null && data !== undefined) {
          this.appsData = data;
        } else {
          this.appsData = [];
        }
        this.apploading = false;
      });
  }

  /* 获取 namespace 列表*/
  queryNamespaces(): void {
      console.log('尝试获取命名空间');
      if (this.cluster === '') {
         return;
      }
      if (this.namespaceloading === true) {
         return;
      }
      this.namespaceloading = true;
      this.mceService.getNamespaces(this.cluster)
        .subscribe((data: any) => {
          if (data !== null && data !== undefined) {
            this.namespaceData = data.items;
          } else {
            this.namespaceData = [];
          }
          this.namespaceloading = false;
        });
  }

  /* 获取 clusters 列表*/
  queryClusters(): void {
    if (this.clustersloading === true) {
      return;
    }
    this.clustersloading = true;
    this.mceService.getClusters()
      .subscribe((data: any) => {
        if (data !== null && data !== undefined) {
          this.clustersData = data.data;
        } else {
          this.clustersData = [];
        }
        this.clustersloading = false;
        console.log(this.clustersData);
      });
  }

  /* 获取 clusters 列表*/
  queryPods(): void {
    if (this.cluster === '') {
      return;
    }

    if (this.namespace === '') {
      return;
    }

    if (this.app === '') {
      return;
    }

    if ( this.podsloading === true) {
      return;
    }
    this.podsloading = true;
    this.mceService.getPods(this.cluster, this.namespace, this.app)
      .subscribe((data: any) => {
        if (data !== null && data !== undefined) {
          this.podsData = [];
          this.total = data.length;
          for (const app of data) {
            let image = '';
            if (app.containers.length > 0) {
              image = app.containers[0].image;
            }
            const info = new PodInfo(app.pod_name, image, app.status, app.start_time);
            this.podsData.push(info);
          }
        } else {
          this.podsData = [];

        }
        this.podsloading = false;
      });
  }

  /*刷新数据*/
  searchData(reset: boolean = false): void {
    if (reset) {

      this.namespaceloading = false;
      this.clustersloading = false;
      this.podsloading = false;
      this.apploading = false;
    }


    // 清理缓存的数据

    this.pageNum = 1;
    this.total = 1;
    this.pageSize = 10;

    this.namespaceData = [];
    this.clustersData = [];
    this.appsData = [];
    this.podsData = [];

    // 获取集群列表
    this.queryClusters();

    // 获取命名空间列表
    this.queryNamespaces();

    // 获取app列表
    this.queryApps();

    // 获取pods列表
    this.queryPods();
  }
}
