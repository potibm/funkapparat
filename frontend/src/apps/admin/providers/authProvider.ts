import { createLogger } from "@core/logger/logger";
import { UserManager } from "oidc-client-ts";
import { AuthProvider } from "react-admin";

const log = createLogger("Auth");

let userManager: UserManager | null = null;
let activeCallbackPromise: Promise<void> | null = null;

export const configureOidc = (authority: string, clientId: string) => {
  userManager = new UserManager({
    authority,
    client_id: clientId,
    redirect_uri: `${globalThis.location.origin}/auth-callback`,
    post_logout_redirect_uri: `${globalThis.location.origin}/`,
    response_type: "code",
    scope: "openid profile email",
  });
};

export const authProvider: AuthProvider = {
  login: async () => {
    if (!userManager) throw new Error("OIDC not configured");
    await userManager.signinRedirect();
  },

  handleCallback: async () => {
    if (!userManager) throw new Error("OIDC not configured");

    if (activeCallbackPromise) {
      return activeCallbackPromise;
    }

    activeCallbackPromise = userManager
      .signinRedirectCallback()
      .then(() => {})
      .catch((error) => {
        log.error("Authentication error in callback:", error);
        throw error instanceof Error ? error : new Error(String(error));
      })
      .finally(() => {
        activeCallbackPromise = null;
      });

    return activeCallbackPromise;
  },

  checkAuth: async () => {
    if (!userManager) return;
    const user = await userManager.getUser();

    if (!user || user.expired) {
      throw new Error("Authentication required");
    }
  },

  logout: async () => {
    if (!userManager) return;

    try {
      const user = await userManager.getUser();
      if (user) {
        await userManager.signoutRedirect();
      } else {
        await userManager.removeUser();
      }
    } catch (error) {
      if (
        error instanceof Error &&
        error.message.includes("No end session endpoint")
      ) {
        log.warn(
          "IdP does not support end session endpoint, clearing local data.",
        );
        await userManager.removeUser();
      } else {
        log.error("Error during logout:", error);
        throw error;
      }
    }
  },

  checkError: async (error: Error | { status?: number }) => {
    if ("status" in error && (error.status === 401 || error.status === 403)) {
      if (userManager) await userManager.removeUser();
      throw new Error("Authentication error");
    }
  },

  getPermissions: async () => {
    return Promise.resolve();
  },
};

export const getAccessToken = async (): Promise<string | null> => {
  if (!userManager) return null;
  const user = await userManager.getUser();
  return user && !user.expired ? user.access_token : null;
};
