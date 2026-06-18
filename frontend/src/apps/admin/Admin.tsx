/* eslint-disable @eslint-react/set-state-in-effect */
import { useState, useRef, useEffect } from "react";
import { Admin, Resource } from "react-admin";
import { BrowserRouter } from "react-router";
import { MyTheme, MyDarkTheme } from "./theme/MyTheme";
import { MyLayout } from "./theme/MyLayout";
import { dataProvider } from "./providers/dataProvider";
import announcements from "./resources/announcement";
import { authProvider, configureOidc } from "./providers/authProvider";
import { OidcLogin } from "./components/auth/OidcLogin";
import { useAppConfig } from "@core/config/useConfig";

export const AdminApp = () => (
  <BrowserRouter>
    <AppBootstrapper />
  </BrowserRouter>
);

export const AppBootstrapper = () => {
  const appConfig = useAppConfig();
  const isOidcActive = appConfig.auth?.type === "oidc";

  const [isConfiguring, setIsConfiguring] = useState(isOidcActive);

  const hasConfiguredRef = useRef(false);

  useEffect(() => {
    if (!isOidcActive) {
      setIsConfiguring(false);
      return;
    }

    if (
      !hasConfiguredRef.current &&
      appConfig.auth?.authority &&
      appConfig.auth?.client_id
    ) {
      configureOidc(appConfig.auth.authority, appConfig.auth.client_id);

      hasConfiguredRef.current = true;

      setIsConfiguring(false);
    }
  }, [isOidcActive, appConfig.auth]);


  if (isConfiguring) {
    return null;
  }

  return (
    <Admin
      authProvider={isOidcActive ? authProvider : undefined}
      loginPage={isOidcActive ? OidcLogin : undefined}
      dataProvider={dataProvider}
      theme={MyTheme}
      darkTheme={MyDarkTheme}
      layout={MyLayout}
      title="Funkapparat Admin"
    >
      <Resource {...announcements} />
    </Admin>
  );
};
