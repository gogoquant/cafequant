import { Injectable } from '@angular/core';
import { ActivatedRoute, PRIMARY_OUTLET } from '@angular/router';
import { Subject } from 'rxjs';
import { IBreadcrumb } from './breadcrumb.model';



@Injectable({
  providedIn: 'root'
})
export class BreadcrumbService {
  private breadcrumbs: IBreadcrumb[] = [];
  private prefixUrl: string;
  breadcrumbSubject: Subject<IBreadcrumb[]>;

  constructor(
    private activatedRoute: ActivatedRoute
  ) {
    this.breadcrumbSubject = new Subject<IBreadcrumb[]>();
  }

  freshBreadcrumbs() {
    const root: ActivatedRoute = this.activatedRoute.root;
    this.breadcrumbs = this.getBreadcrumbs(root, this.prefixUrl);
    this.breadcrumbSubject.next(this.breadcrumbs);
  }

  setPrefixUrl(prefixUrl: string) {
    this.prefixUrl = prefixUrl;
  }

  /**
   * Returns array of IBreadcrumb objects that represent the breadcrumb
   *
   * @class DetailComponent
   * @method getBreadcrumbs
   * @param {ActivateRoute} route
   * @param {string} url
   * @param {IBreadcrumb[]} breadcrumbs
   */
  private getBreadcrumbs(route: ActivatedRoute, url: string = '', breadcrumbs: IBreadcrumb[] = []): IBreadcrumb[] {
    const ROUTE_DATA_BREADCRUMB = 'breadcrumb';

    // get the child routes
    const children: ActivatedRoute[] = route.children;

    // return if there are no more children
    if (children.length === 0) {
      return breadcrumbs;
    }

    // iterate over each children
    for (const child of children) {
      // verify primary route
      if (child.outlet !== PRIMARY_OUTLET) {
        continue;
      }

      // verify the custom data property "breadcrumb" is specified on the route
      if (!child.snapshot.data.hasOwnProperty(ROUTE_DATA_BREADCRUMB)) {
        return this.getBreadcrumbs(child, url, breadcrumbs);
      }

      // get the route's URL segment
      const routeURL: string = child.snapshot.url.map(segment => segment.path).join('/');

      // append route URL to URL
      url += `/${routeURL}`;

      // add breadcrumb
      const breadcrumb: IBreadcrumb = {
        label: child.snapshot.data[ROUTE_DATA_BREADCRUMB],
        params: child.snapshot.params,
        url: url
      };
      console.log('label=', child.snapshot.data[ROUTE_DATA_BREADCRUMB]);
      console.log('url=', url);
      breadcrumbs.push(breadcrumb);

      // recursive
      return this.getBreadcrumbs(child, url, breadcrumbs);
    }
  }
}
