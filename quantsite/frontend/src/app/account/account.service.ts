import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders, HttpParams } from '@angular/common/http';
import { LoggerService } from './../commons/services/logger.service';
import { Observable, of } from 'rxjs';
import { catchError, map, tap } from 'rxjs/operators';

const httpOptions = {
  headers: new HttpHeaders({ 'Content-Type': 'application/json' })
};

@Injectable({
  providedIn: 'root'
})
export class AccountService {

  private accountsUrl = '/server/account/api';  // URL to web api

  constructor(private http: HttpClient, private logger: LoggerService) { }

  /** GET accounts from the server */
  getAccounts(keywords: string = '', limit: string = '', pageIndex: number = 1, pageSize: number = 10): Observable<{}> {
    const params = new HttpParams()
      .append('keywords', `${keywords}`)
      .append('limit', `${limit}`)
      .append('page_index', `${pageIndex}`)
      .append('page_size', `${pageSize}`);
    return this.http.get(`${this.accountsUrl}/users`, {
      params
    });
  }


  syncUsers(): Observable<{}> {
    const options = { headers: new HttpHeaders({ 'Content-Type': 'application/json' }) };

    this.logger.info('sync user');
    return this.http.post(`${this.accountsUrl}/users/sync`, JSON.stringify({
      'Keyword': '*'
    }), options);
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

  /** GET userId from the server */
  getUserIdByUserName(userName: string = ''): Observable<{}> {
    return this.http.get(`${this.accountsUrl}/users/${userName}`);
  }
}
