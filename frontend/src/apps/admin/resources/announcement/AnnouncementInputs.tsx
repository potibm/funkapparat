import { TextInput, BooleanInput, required } from "react-admin";
import { MarkdownInput } from "@admin/components/inputs/MarkdownInput";

export const AnnouncementInputs = () => {
  return (
    <>
      <TextInput source="title" required />

      <MarkdownInput source="body" label="Body" validate={[required()]} />

      <BooleanInput source="is_hidden" label="Hidden" />

      <BooleanInput source="is_urgent" label="Urgent" />

      <TextInput source="external_url" label="External URL" />
    </>
  );
};
