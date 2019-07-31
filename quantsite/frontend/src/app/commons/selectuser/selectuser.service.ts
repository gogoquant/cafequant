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
export class SelectuserService {

  private accountsUrl = '/server/account/api';  // URL to web api

  constructor(private http: HttpClient, private logger: LoggerService) { }

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
}
