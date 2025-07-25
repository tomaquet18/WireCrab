import React from 'react'
import {createRoot} from 'react-dom/client'
import './style.css'
import App from './App'
import { BrowserRouter, Routes, Route } from 'react-router-dom'
import WelcomePage from './pages/WelcomePage'

const container = document.getElementById('root')

const root = createRoot(container!)

root.render(
    <React.StrictMode>
        <BrowserRouter>
            <Routes>
                <Route path="/" element={<App />} />
                <Route path="/welcome" element={<WelcomePage />} />
            </Routes>
        </BrowserRouter>
    </React.StrictMode>
)
