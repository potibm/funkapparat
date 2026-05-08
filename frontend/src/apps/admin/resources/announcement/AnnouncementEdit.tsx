import { Edit, SimpleForm } from "react-admin";
import { AnnouncementInputs } from "./AnnouncementInputs";

export const AnnouncementEdit = () => {
  return (
    <Edit title="Edit Annoucement">
      <SimpleForm>
        <AnnouncementInputs />
      </SimpleForm>
    </Edit>
  );
};

export default AnnouncementEdit;
