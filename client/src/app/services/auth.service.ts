import {Injectable, signal, WritableSignal} from '@angular/core';
import {HttpClient} from "@angular/common/http";

interface User {
	name: string;
	password: string;
}

@Injectable({
	providedIn: 'root'
})
export class AuthService {
	loginFlowStatus: WritableSignal<boolean> = signal(false);
	private url: string = "http://localhost:8080/client";

	constructor(private http: HttpClient) {
	}

	auth(username: string, password: string) {
		const user: User = {name: username, password: password};
		return this.http.post(this.url + "/login", user)
			.subscribe({
				complete: () => {
					this.http.get<void>(this.url + "/redirect");
				}
			});
	}

	register(username: string, password: string) {
		const user: User = {name: username, password: password};
		return this.http.post<void>(this.url + "/register", user);
	}

	getLoginFlowStatus() {
		return this.loginFlowStatus();
	}

	setLoginFlowStatus(value: boolean) {
		return this.loginFlowStatus.set(value)
	}
}
