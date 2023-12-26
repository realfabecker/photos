export type ResponseDTO<T> = {
  status: "success" | "error";
  data: T;
};

export type PagedDTO<T> = {
  page_count: number;
  items: T[];
  page_token?: string;
  has_more: boolean;
};

export type Photo = {
  id: string;
  tags: string[];
  url: string;
  fileName?: string;
  title: string;
  createdAt: string;
};

export enum ActionStatus {
  IDLE = "idle",
  DONE = "done",
  ERROR = "error",
  LOADING = "loading",
}

export interface LoginDTO {
  RefreshToken: string;
  AccessToken: string;
}

export enum RoutesEnum {
  Login = "/login",
  Photos = "/photos",
}

export enum ProviderEnum {
  Lambda = "lambda",
  Picsum = "picsum",
}

export enum ModalState {
  Open = "open",
  Closed = "closed",
}
