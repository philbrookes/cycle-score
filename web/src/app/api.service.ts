import { Injectable } from '@angular/core';
import { Headers, Http, RequestOptions, Response, URLSearchParams } from '@angular/http';
import { environment } from '../environments/environment'
import 'rxjs/add/operator/toPromise';
import { Observable } from 'rxjs/Observable';
import { Subject } from 'rxjs/Subject';
import { Observer } from "rxjs/Observer";


@Injectable()
export class ApiService {
  constructor(private http: Http) { }

  public isAuthorized():Promise<any> {
    return this.get(environment.apiHost + "/api/auth/check")
  }

  public oauthUrl():Promise<any> {
    return this.get(environment.apiHost + "/api/auth/url")
  }

  public getScore(from, to):Promise<any> {
    return this.get(environment.apiHost + "/api/score/generate/from/"+ from + "/to/" +to)
  }

  private get(url):Promise<any>{
    return this.http.get(url).toPromise()
      .then(response => {
        return response.json();
      })
      .catch((err) => {
        console.log("problem getting URL: " + url);
        console.log(err)
      })
  }

}
