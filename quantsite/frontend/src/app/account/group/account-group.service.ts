import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders, HttpParams } from '@angular/common/http';
import { LoggerService } from './../../commons/services/logger.service';
import { Observable, of } from 'rxjs';

const httpOptions = {
  headers: new HttpHeaders({ 'Content-Type': 'application/json' })
};

@Injectable({
  providedIn: 'root'
})
export class AccountGroupService {

  private accountsUrl = '/server/account/api';  // URL to web api

  constructor(private http: HttpClient, private logger: LoggerService) { }

  /** GET modules from the server */
  getModules(): Observable<{}> {
    return this.http.get(`${this.accountsUrl}/modules`);
  }

  /** GET accounts from the server */
  getGroupsById(id: number): Observable<{}> {
    const params = new HttpParams()
      .append('project_id', `${id}`);
    return this.http.get(`${this.accountsUrl}/groups/{project_id}/findByProjectId`, {
      params
    });
  }

  /** 添加用户组 */
  AddGroup(params): Observable<{}> {
    return this.http.post(`${this.accountsUrl}/groups`,
      params
    );
  }

  /** 删除用户组 */
  DelGroup(id): Observable<{}> {
    return this.http.delete(`${this.accountsUrl}/groups/${id}`);
  }

  /** 编辑用户组 */
  EditGroup(params): Observable<{}> {
    return this.http.patch(`${this.accountsUrl}/groups`,
      params
    );
  }

  /** 根据用户组ID获取用户 */
  getUserByGroup(id): Observable<{}> {
    const params = new HttpParams()
      .append('group_id', `${id}`);
    return this.http.get(`${this.accountsUrl}/users_groups/{group_id}/findByGroupId`, {
      params
    });
  }

  /** 删除用户 */
  DelUser(id): Observable<{}> {
    return this.http.delete(`${this.accountsUrl}/users_groups/${id}`);
  }

  /** 添加用户 */
  AddUsers(params): Observable<{}> {
    return this.http.post(`${this.accountsUrl}/users_groups`,
      params
    );
  }

   /** GET accounts from the server */
   getAccounts(keywords: string = '', limit: string = '', pageIndex: number = 1, pageSize: number = 20): Observable<{}> {
    const params = new HttpParams()
      .append('keywords', `${keywords}`)
      .append('limit', `${limit}`)
      .append('page_index', `${pageIndex}`)
      .append('page_size', `${pageSize}`);
    return this.http.get(`${this.accountsUrl}/users`, {
      params
    });
  }

  /**
   * Handle Http operation that failed.
   * Let the app continue.
   * @param operation - name of the operation that failed
   * @param result - optional value to return as the observable result
   */
  private handleError<T>(operation = 'operation', result ?: T) {
    return (error: any): Observable<T> => {
      // TODO: send the error to remote logging infrastructure
      this.logger.error(error); // log to console instead
      // Let the app keep running by returning an empty result.
      return of(result as T);
    };
  }
}
