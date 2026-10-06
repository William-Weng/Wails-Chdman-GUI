import { mount } from 'svelte'
import App from './App.svelte'
import "./less/reset.less"
import "./less/main.less"

mount(App, { target: document.getElementById('app')! })
