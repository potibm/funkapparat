import { StrictMode } from "react";
import { createRoot } from "react-dom/client";

import "./index.css";
import App from "./App.tsx";

import { ConfigContext } from "@core/config/ConfigContext.tsx";
import { createLogger } from "@core/logger/logger.ts";
import { AppConfigSchema } from "@core/config/config.schemas.ts";
import { configureOidc } from "@admin/providers/authProvider.ts";
import * as Sentry from "@sentry/react";

const log = createLogger("Bootstrapper");
const API_HOST = import.meta.env.VITE_API_HOST ?? "http://localhost:3101";

export async function bootstrapApp() {
  const rootElement = document.getElementById("root");
  if (!rootElement) throw new Error("Failed to find the root element");
  const root = createRoot(rootElement);

  try {
    // 1. Fetch config and validate
    const controller = new AbortController();
    const timeoutId = globalThis.setTimeout(() => controller.abort(), 10000);

    const res = await fetch(`${API_HOST}/api/config`, {
      signal: controller.signal,
    });
    if (!res.ok) throw new Error(`Config error: ${res.statusText}`);

    const data = await res.json();
    globalThis.clearTimeout(timeoutId);

    const config = AppConfigSchema.parse(data);
    log.debug("Config loaded:", config);

    // 2. Initialize Sentry
    if (config.sentry?.dsn && !Sentry.isInitialized()) {
      log.debug("Configuring Sentry");
      Sentry.init({
        dsn: config.sentry.dsn,
        environment: config.sentry.environment,
        release: config.sentry.version,
        replaysSessionSampleRate: config.sentry.replay_session_sample_rate,
        replaysOnErrorSampleRate: config.sentry.replay_error_sample_rate,
        integrations: [
          Sentry.replayIntegration(),
          Sentry.browserTracingIntegration(),
        ],
      });
    }

    // 3. Initialize OIDC
    if (
      config.auth?.type === "oidc" &&
      config.auth.authority &&
      config.auth.client_id
    ) {
      log.debug("Configuring OIDC");
      configureOidc(config.auth.authority, config.auth.client_id);
    }

    // 4. Start React
    root.render(
      <StrictMode>
        <ConfigContext value={config}>
          <App />
        </ConfigContext>
      </StrictMode>,
    );
  } catch (err) {
    log.error("Bootstrap failed:", err);

    root.render(
      <div style={{ padding: 20, color: "red", fontFamily: "sans-serif" }}>
        <h2>System Configuration Error</h2>
        <pre>{err instanceof Error ? err.message : "Unknown error"}</pre>
      </div>,
    );
  }
}

if (!import.meta.env.TEST) {
  bootstrapApp();
}
