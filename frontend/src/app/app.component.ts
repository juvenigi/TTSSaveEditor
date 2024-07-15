import {Component} from '@angular/core';
import {RouterOutlet} from '@angular/router';
import {AsyncPipe, NgComponentOutlet} from "@angular/common";
import {HeaderComponent} from "./components/header/header.component";
import {DirectoryPageComponent} from "./pages/directory/directory.page.component";

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [RouterOutlet, AsyncPipe, HeaderComponent, NgComponentOutlet, DirectoryPageComponent],
  providers: [],
  templateUrl: './app.component.html',
  styleUrl: './app.component.css'
})
export class AppComponent {
}
