import { Component } from '@angular/core';
import {InputText} from "primeng/inputtext";
import {FloatLabel} from "primeng/floatlabel";
import {FormBuilder, FormGroup, FormsModule, ReactiveFormsModule, Validators} from "@angular/forms";
import {Password} from "primeng/password";
import {Button} from "primeng/button";
import {Router, RouterLink} from "@angular/router";
import {AuthService} from "../../services/auth.service";

@Component({
  selector: 'app-login',
	imports: [
		InputText,
		FloatLabel,
		FormsModule,
		Password,
		ReactiveFormsModule,
		Button,
		RouterLink
	],
  templateUrl: './login.component.html',
  styleUrl: './login.component.css'
})

export class LoginComponent {
	form!: FormGroup;

	constructor(private formBuilder: FormBuilder,
				private authService: AuthService,
				private router: Router
	) {
		this.authService.setLoginFlowStatus(true);
		const v = Validators;
		this.form = this.formBuilder.group({
			username: ['', v.required],
			password: ['', v.required]
		})
	}

	onSubmit() {
		this.authService.auth(this.form.value.username, this.form.value.password);
	}
}
