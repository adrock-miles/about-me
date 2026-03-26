import React from "react";
import ReactDOM from "react-dom/client";
import { BrowserRouter } from "react-router-dom";
import { Theme } from "@radix-ui/themes";
import "@radix-ui/themes/styles.css";
import App from "./App";
import "./styles/global.css";

ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <BrowserRouter>
      <Theme appearance="dark" accentColor="gray" radius="large" scaling="100%">
        <App />
      </Theme>
    </BrowserRouter>
  </React.StrictMode>,
);
