import { Injectable } from '@angular/core';
import { HttpHeaders, HttpParams, HttpClient } from '@angular/common/http';
import { Observable, of } from 'rxjs';
import { LoggerService } from '../../commons/services/logger.service';
import { Base64 } from 'js-base64';


@Injectable({
  providedIn: 'root'
})
export class AgileService {

  private agileUrl = '/server/agile-admin/v1';

  constructor(private http: HttpClient,
    private logger: LoggerService) { }


  /** GET products from the server */
  getProducts(pageNum: number, pageSize: number, username: string, creator: string, productName: string, tenantId: string): Observable<{}> {
    const base64String = Base64.encode(username);
    const headers = new HttpHeaders({
      'Content-Type': 'application/json',
      'X-Matrix-Username': username,
      'X-Matrix-Token': base64String
    });

    const params = new HttpParams()
      .append('pageNum', `${pageNum}`)
      .append('pageSize', `${pageSize}`)
      .append('creator', `${creator}`)
      .append('productName', `${productName}`)
      .append('tenantId', `${tenantId}`);

    // this.logger.info(params);
    return this.http.get(`${this.agileUrl}/products`,
      {
        headers,
        params
      });
  }

  /** GET projects from the server */
  getProjects(pageNum: number, pageSize: number, username: string, creator: string,
    productName: string, projectName: string, tenantId: string): Observable<{}> {
    const base64String = Base64.encode(username);
    const headers = new HttpHeaders({
      'Content-Type': 'application/json',
      'X-Matrix-Username': username,
      'X-Matrix-Token': base64String
    });

    const params = new HttpParams()
      .append('pageNum', `${pageNum}`)
      .append('pageSize', `${pageSize}`)
      .append('creator', `${creator}`)
      .append('productName', `${productName}`)
      .append('projectName', `${projectName}`)
      .append('tenantId', `${tenantId}`);

    // this.logger.info(params);
    return this.http.get(`${this.agileUrl}/projects`,
      {
        headers,
        params
      });
  }

  /** GET project from the server */
  getProjectById(projectId: string, username: string): Observable<{}> {
    const base64String = Base64.encode(username);
    const headers = new HttpHeaders({
      'Content-Type': 'application/json',
      'X-Matrix-Username': username,
      'X-Matrix-Token': base64String
    });

    return this.http.get(`${this.agileUrl}/projects/${projectId}`,
      {
        headers
      });
  }


  /** get Project configs */
  getProjectConfig(projectId: string, username: string): Observable<{}> {
    const base64String = Base64.encode(username);
    const headers = new HttpHeaders({
      'Content-Type': 'application/json',
      'X-Matrix-Username': username,
      'X-Matrix-Token': base64String
    });

    return this.http.get(`${this.agileUrl}/projectconfs/${projectId}`, { headers });
  }


  /** get Project configs */
  getProjectUsers(projectId: string, username: string): Observable<{}> {
    const base64String = Base64.encode(username);
    const headers = new HttpHeaders({
      'Content-Type': 'application/json',
      'X-Matrix-Username': username,
      'X-Matrix-Token': base64String
    });
    const params = new HttpParams().append('projectId', `${projectId}`);
    return this.http.get(`${this.agileUrl}/projects/user-role-relations`, { headers, params });
  }

  /** get Roles */
  getRoles(username: string): Observable<{}> {
    const base64String = Base64.encode(username);
    const headers = new HttpHeaders({
      'Content-Type': 'application/json',
      'X-Matrix-Username': username,
      'X-Matrix-Token': base64String
    });
    return this.http.get(`${this.agileUrl}/roles`, { headers });
  }


  /** GET projects from the server test */
  gettests(): Observable<{}> {
    return this.http.get(`${this.agileUrl}/projects/test`);
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
