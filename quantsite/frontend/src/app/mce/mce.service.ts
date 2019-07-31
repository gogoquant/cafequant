import { Injectable } from '@angular/core';
import { HttpHeaders, HttpParams, HttpClient } from '@angular/common/http';
import { Observable, of } from 'rxjs';
import { LoggerService } from '../commons/services/logger.service';
import { Base64 } from 'js-base64';

@Injectable({
  providedIn: 'root'
})

export class MceService {

  // mce api address
  private mceUrl = '/server/mce-admin/v1';
  private token = 'eyJhbGciOiJSUzI1NiIsImtpZCI6IjcxNTYxNmEzNDU0OTcwZjBjMzA0YmIyN2E2N2E1YWYxOThlODY1MWMifQ.eyJpc3MiOiJodHRwczovL2FjY291bnQuY2hhbmdob25nLmlvIiwic3ViIjoiQ2l4MWFXUTlhRzl1WjNkbGFTNXRaV2tzYjNVOWRYTmxjbk1zWkdNOVkyaGhibWRvYjI1bkxHUmpQV052YlJJRWJHUmhjQSIsImF1ZCI6ImV4YW1wbGUtYXBwIiwiZXhwIjoxNTM1NTM3NDAyLCJpYXQiOjE1MzU0NTEwMDIsImF0X2hhc2giOiJ4MzJrUXZCMGtOc1V6TW5uc0RWYmVnIiwiZW1haWwiOiJob25nd2VpLm1laUBjaGFuZ2hvbmcuY29tIiwiZW1haWxfdmVyaWZpZWQiOnRydWUsIm5hbWUiOiLmooXnuqLkvJ8ifQ.MDqCIZ6yxJafxRP12WmlmkqZ2STm3Chp7SZ0W1ypO7dy237P12fvx9KGZOldOHl21bAuxNVjKESHvp5Oy-GF9N7cp29_RQxR48VjOhuHrv2vOuUeXkztFLT5eWzwDi9fcyRrKjx88M_egNuLBifCe2yGJntA6bckPJUe53T_BJxu1nVv2SXnCVlbfhVc03rdt9TCWXwuC46ZUC6fW8s8UVL95U9iOqxB_nfuM17rxALBxUZIUFTRIOgNfMicAaNOjm5EIm019sKWEJrikPpjcYCdymoj3K-K97t5olCIrETGhfJ8T8A3hGf_jGQhS47E2UiuPuFl1xOl609YP1NerQ';

  constructor(private http: HttpClient,
              private logger: LoggerService) { }



  /* get namespace from the server */
  getNamespaces(cluster_id: string): Observable<{}> {
    const headers = new HttpHeaders({
      'Content-Type': 'application/json',
      'X-Matrix-Token': `${this.token}`,
    });

    const params = new HttpParams();

    this.logger.info(params);
    return this.http.get(`${this.mceUrl}/clusters/${cluster_id}/namespaces`,
      {
        headers,
        params
      });
  }

  /* get app from the server */
  getApps(cluster_id: string, namespace: string): Observable<{}> {
    const headers = new HttpHeaders({
      'Content-Type': 'application/json',
      'X-Matrix-Token': `${this.token}`,
    });

    const params = new HttpParams();

    this.logger.info(params);
    return this.http.get(`${this.mceUrl}/clusters/${cluster_id}/namespaces/${namespace}/applicationgroups`,
      {
        headers,
        params
      });
  }

  /* get clusters from the server */
  getClusters(): Observable<{}> {
    const headers = new HttpHeaders({
      'Content-Type': 'application/json',
      'X-Matrix-Token': `${this.token}`,
    });

    const params = new HttpParams();

    this.logger.info(params);
    return this.http.get(`${this.mceUrl}/clusters`,
      {
        headers,
        params
      });
  }

  /* get pods from the server */
  getPods(cluster_id: string, namespace: string, app: string): Observable<{}> {
    // const base64app = Base64.encode(app);
    const headers = new HttpHeaders({
      'Content-Type': 'application/json',
      'X-Matrix-Token': `${this.token}`,
    });

    const params = new HttpParams()
      .append('app-name', app);

    this.logger.info(params);
    return this.http.get(`${this.mceUrl}/clusters/${cluster_id}/namespaces/${namespace}/pods`,
      {
        headers,
        params
      });
  }

  /** GET products from the server */
  gettests(): Observable<{}> {
    return this.http.get(`${this.mceUrl}/projects/test`);
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
