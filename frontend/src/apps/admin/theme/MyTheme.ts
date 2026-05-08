import { defaultTheme } from "react-admin";
import type { ThemeOptions } from "@mui/material";

export const MyTheme: ThemeOptions = {
  ...defaultTheme,
  palette: {
    ...defaultTheme.palette,
    mode: "light",
    primary: {
      // Ein sehr dunkles Schiefergrau für perfekte Lesbarkeit auf hellem Grund
      main: "#1E293B",
      contrastText: "#FFFFFF",
    },
    secondary: {
      // Ein klares, technisches Blau als Sekundärfarbe
      main: "#3B82F6",
      contrastText: "#FFFFFF",
    },
    background: {
      default: "#F8FAFC",
      paper: "#FFFFFF",
    },
    error: { main: "#FF3366" },
    warning: { main: "#FF9800" },
    info: { main: "#00BCD4" },
    success: { main: "#10B981" },
  },
  components: {
    ...defaultTheme.components,
    // HIER SETZEN WIR DEN AKZENT FÜR DEN LIGHT MODE
    MuiAppBar: {
      styleOverrides: {
        root: {
          backgroundColor: "#CCFF00", // Der knallige Limetten-Header
          color: "#000000", // Schwarzer Text/Icons im Header für den Kontrast
        },
      },
    },
    MuiButton: {
      styleOverrides: {
        root: {
          borderRadius: "6px",
          fontWeight: 600,
          textTransform: "none",
        },
      },
    },
  },
};

export const MyDarkTheme: ThemeOptions = {
  ...defaultTheme,
  palette: {
    ...defaultTheme.palette,
    mode: "dark",
    primary: {
      main: "#CCFF00", // Limette als Hauptfarbe im Dark Mode
      light: "#D9FF33",
      dark: "#99CC00",
      contrastText: "#000000", // WICHTIG: Schwarzer Text auf den Limetten-Buttons
    },
    secondary: {
      // Violett ist weg. Wir nutzen ein strahlendes Cyan als Kontrast zur Limette
      main: "#00E5FF",
      contrastText: "#000000",
    },
    background: {
      default: "#0F172A",
      paper: "#1E293B",
    },
    error: { main: "#FF3366" },
    warning: { main: "#FFB300" },
    info: { main: "#00E5FF" },
    success: { main: "#00E676" },
  },
  components: {
    ...defaultTheme.components,
    MuiButton: {
      styleOverrides: {
        root: {
          borderRadius: "6px",
          fontWeight: 600,
          textTransform: "none",
        },
      },
    },
    MuiTableCell: {
      styleOverrides: {
        head: {
          backgroundColor: "#0F172A",
        },
      },
    },
    MuiChip: {
      styleOverrides: {
        filledPrimary: {
          backgroundColor: "#CCFF00",
          color: "#000000",
        },
        root: {
          backgroundColor: "rgba(15, 23, 42, 0.6)",
          border: "1px solid rgba(204, 255, 0, 0.3)", // Dezent limettengrüner Rand
        },
        label: {
          color: "#CCFF00", // Limettengrüner Text
        },
      },
    },
    MuiTableRow: {
      styleOverrides: {
        root: {
          "&:hover": {
            backgroundColor: "rgba(255, 255, 255, 0.04) !important",
          },
        },
      },
    },
  },
};
