import {Component} from '@angular/core';
import {Button} from "primeng/button";
import {FloatLabel} from "primeng/floatlabel";
import {FormBuilder, FormGroup, FormsModule, ReactiveFormsModule, Validators} from "@angular/forms";
import {InputText} from "primeng/inputtext";
import {Password} from "primeng/password";
import {Router} from "@angular/router";
import {AuthService} from "../../services/auth.service";
import {Toast, ToastModule} from "primeng/toast";
import {MessageService} from "primeng/api";

@Component({
	selector: 'app-register',
	imports: [
		Button,
		FloatLabel,
		FormsModule,
		InputText,
		Password,
		ReactiveFormsModule,
		Toast,
		ToastModule,
	],
	providers: [MessageService],
	templateUrl: './register.component.html',
	styleUrl: './register.component.css'
})
export class RegisterComponent {
	form!: FormGroup;

	constructor(private formBuilder: FormBuilder,
				private authService: AuthService,
				private router: Router,
				private messageService: MessageService
	) {
		this.authService.setLoginFlowStatus(true);
		const v = Validators;
		this.form = this.formBuilder.group({
			username: ['', v.required],
			password: ['', v.required]
		})
	}

	onSubmit() {
		this.authService.register(this.form.value.username, this.form.value.password)
			.subscribe({
				complete: () => {
					this.router.navigate(['/login']);
				},
				error: (err) => {
					this.messageService.add({
						severity: 'error',
						summary: 'Error',
						detail: 'Failed to create user. Please try again',
						key: "registration-error",
						life: 2000
					});
				}
			});
	}
}
