import { Injectable } from '@angular/core';
import { HttpHeaders, HttpClient, HttpParams } from '@angular/common/http';
import { LoggerService } from 'src/app/commons/services/logger.service';
import { Observable, of } from 'rxjs';
import { Base64 } from 'js-base64';

const httpOptions = {
  headers: new HttpHeaders({ 'Content-Type': 'application/json' })
};



@Injectable({
  providedIn: 'root'
})
export class TenantService {

  private accountsUrl = '/server/account/api';  // URL to web api

  constructor(private http: HttpClient, private logger: LoggerService) { }

  /** 添加租户 */
  AddTenant(params): Observable<{}> {
    return this.http.post(`${this.accountsUrl}/tenants`,
      params
    );
  }


  /** GET tenants from the server */
  queryTenants(pageNum: number, pageSize: number, keyword: string): Observable<{}> {
    const headers = new HttpHeaders({
      'Content-Type': 'application/json',
    });

    const params = new HttpParams()
      .append('page_index', `${pageNum}`)
      .append('page_size', `${pageSize}`)
      .append('keyword', `${keyword}`);
    return this.http.get(`${this.accountsUrl}/tenants`,
      {
        headers,
        params
      });
  }

  /** delete tenant by tenant_uuid */
  deleteTenant(tenant_uuid: string): Observable<{}> {
    const headers = new HttpHeaders({
      'Content-Type': 'application/json',
    });
    return this.http.delete(`${this.accountsUrl}/tenants/${tenant_uuid}`,
      {
        headers
      });
  }

  /**
   * Handle Http operation that failed.
   * Let the app continue.
   * @param operation - name of the operation that failed
   * @param result - optional value to return as the observable result
   */
  private handleError<T>(operation = 'operation', result?: T) {
    return (error: any): Observable<T> => {
      // TODO: send the error to remote logging infrastructure
      this.logger.error(error); // log to console instead
      // Let the app keep running by returning an empty result.
      return of(result as T);
    };
  }
}
