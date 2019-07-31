import { Injectable } from '@angular/core';
import { HttpClient, HttpParams, HttpErrorResponse, HttpHeaders } from '@angular/common/http';
import { Observable, of, throwError  } from 'rxjs';
import { catchError } from 'rxjs/operators';

import { AppInfo, Organization } from './beans';
import { APP_INFO_LIST, DEPARTMENTS } from './mock-data';

@Injectable()
export class MseService {
  private urlPrefix = '/server/microservice/api';
  private appInfoCache: AppInfo[];

  constructor(private httpClient: HttpClient) {
    this.appInfoCache = [];
   }

  queryApps(orgName: string, owner: string, page: number = 0, size: number = 10): Observable<AppInfo[]> {
    const options = {
      headers: new HttpHeaders()
        .append('Authorization', 'Basic c2h1aG9uZy5saTphZG1pbg=='),
      params: new HttpParams()
        .append('orgName', orgName ? orgName : '')
        .append('owner', owner ? owner : '')
        .append('page', page.toString())
        .append('size', size.toString())
    };
    return this.httpClient.get<AppInfo[]>(`${this.urlPrefix}/apps/find`, options)
      .pipe(
        catchError(this.handleError)
      );
  }

  queryOrganizations(): Observable<Organization[]> {
    const options = {
      headers: new HttpHeaders()
        .append('Authorization', 'Basic c2h1aG9uZy5saTphZG1pbg==')
    };
    return this.httpClient.get<AppInfo[]>(`${this.urlPrefix}/organizations`, options)
      .pipe(
        catchError(this.handleError)
      );
  }

  queryAppInfo(id: number): Observable<AppInfo> {
    return of<AppInfo>(this.appInfoCache.find(appInfo => appInfo.id === id));
  }

  private handleError(error: HttpErrorResponse) {
    if (error.error instanceof ErrorEvent) {
      // A client-side or network error occurred. Handle it accordingly.
      console.error('An error occurred:', error.error.message);
    } else {
      // The backend returned an unsuccessful response code.
      // The response body may contain clues as to what went wrong,
      console.error(
        `Backend returned code ${error.status}, ` +
        `body was: ${error.error}`);
    }
    // return an observable with a user-facing error message
    return throwError(
      'Something bad happened; please try again later.');
  }
}
