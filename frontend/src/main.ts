import { mount } from 'svelte'
import App from './App.svelte'
import "./main.less"

mount(App, { target: document.getElementById('app')! })
