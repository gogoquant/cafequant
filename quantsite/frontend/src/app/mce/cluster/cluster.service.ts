import { Injectable } from "@angular/core";
import { HttpHeaders, HttpClient, HttpParams } from "@angular/common/http";
import { Observable } from "rxjs";

const httpOptions = {
  headers: new HttpHeaders({ "Content-Type": "application/json" })
};

@Injectable({
  providedIn: "root"
})
export class ClusterService {
  private clusterUrl = "/server/mce-admin/v1";
  constructor(private http: HttpClient) {}

  //调用查询集群接口
  getClusters(pageIndex: number = 1, pageSize: number = 10): Observable<{}> {
    const params = new HttpParams()
      .append("page", `${pageIndex}`)
      .append("page_size", `${pageSize}`);
    return this.http.get(`${this.clusterUrl}/clusters`, { params });
  }
}
