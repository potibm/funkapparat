import { RaRecord } from "react-admin";

export interface AnnouncementRecord extends RaRecord {
  id: string | number;
  title: string;
  body: string;
  body_html: string;
  body_plain: string;
  is_hidden: boolean;
  is_urgent: boolean;
  external_url?: string;
}
