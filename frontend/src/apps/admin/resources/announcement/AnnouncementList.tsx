import {
  List,
  Datagrid,
  DateField,
  TextField,
  SearchInput,
  FunctionField,
  TopToolbar,
  CreateButton,
  ExportButton,
} from "react-admin";
import { Chip } from "@mui/material";
import { AnnouncementRecord } from "./announcements.types";
import { BooleanToggleField } from "@admin/components/fields/BooleanToggleField";
import { useAppConfig } from "@core/config/useConfig";

const announcementFilters = [
  <SearchInput key="q" source="q" alwaysOn placeholder="Search by title..." />,
];

const ListActions = () => (
  <TopToolbar>
    <CreateButton />
    <ExportButton />
  </TopToolbar>
);

export const ScheduleEntriesList = () => {
  const { date_locale: dateLocale, date_options: dateOptions } = useAppConfig();

  return (
    <List
      title="Annoucements"
      sort={{ field: "id", order: "DESC" }}
      actions={<ListActions />}
      filters={announcementFilters}
    >
      <Datagrid rowClick="edit" bulkActionButtons={false}>
        {/* 1. STATUS FIELD */}
        <FunctionField
          label="Status"
          sortable={false}
          render={(record: AnnouncementRecord) => {
            if (record.is_hidden) {
              return (
                <Chip
                  label="Hidden"
                  color="success"
                  size="small"
                  variant="filled"
                />
              );
            } else if (record.is_urgent) {
              return <Chip label="Urgent" size="small" variant="outlined" />;
            } else {
              return (
                <Chip
                  label="Visible"
                  color="primary"
                  size="small"
                  variant="outlined"
                />
              );
            }
          }}
        />

        {/* 2. ID */}
        <TextField source="id" label="ID" />

        {/* 3. CREATED AT */}
        <DateField
          source="created_at"
          showTime={true}
          emptyText="-"
          locales={dateLocale}
          options={dateOptions}
        />

        {/* 4. TITLE */}
        <TextField source="title" label="Title" />

        {/* 5. BODY */}
        <TextField
          source="body_plain"
          label="Body"
          sx={{
            display: "block",
            maxWidth: "250px",
            whiteSpace: "nowrap",
            overflow: "hidden",
            textOverflow: "ellipsis",
          }}
        />

        {/* 6. HIDDEN TOGGLE */}
        <BooleanToggleField source="is_hidden" label="Hidden" />
      </Datagrid>
    </List>
  );
};

export default ScheduleEntriesList;
