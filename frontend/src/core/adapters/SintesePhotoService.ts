import "reflect-metadata";
import { inject, injectable } from "inversify";
import { type IAuthService, IPhotoService, Types } from "@core/ports/ports.ts";
import { PagedDTO, Photo, ResponseDTO } from "@core/domain/domain.ts";

type SinteseResDTO<T = Record<string, any>> = {
  status: "success" | "error";
  data: T;
};

type SintesePagedDTO<T> = {
  page_count: number;
  items: T[];
  has_more: boolean;
  page_token?: string;
};

type SintesePhoto = {
  photoId: string;
  userId: string;
  fileName: string;
  title: string;
  url: string;
  createdAt: string;
};

@injectable()
export class SintesePhotoService implements IPhotoService {
  constructor(
    @inject(Types.AuthService) private readonly auth: IAuthService,
    private readonly baseUrl: string = import.meta.env.VITE_API_BASE_URL
  ) {}
  async fetchPhotos(opts: {
    page: number;
    limit: number;
    token?: string;
  }): Promise<ResponseDTO<PagedDTO<Photo>>> {
    const params = new URLSearchParams();
    params.set("limit", opts.limit + "");
    params.set("created_at", "2023");
    if (opts.token) params.set("page_token", opts.token + "");
    const res = await fetch(`${this.baseUrl}/midia?${params.toString()}`, {
      headers: { Authorization: `Bearer ${this.auth.getAccessToken()}` },
    });
    const data = (await res.json()) as SinteseResDTO<
      SintesePagedDTO<SintesePhoto>
    >;
    const items: Photo[] = data.data.items.map((p) => ({
      url: p.url,
      id: p.photoId,
      tags: ["tag"],
      title: p.title,
      createdAt: p.createdAt,
    }));
    return {
      status: "success",
      data: {
        items: items,
        page_count: data.data.page_count,
        page_token: data.data.page_token,
        has_more: data.data.has_more,
      },
    };
  }

  async getUploadUrl(name: string): Promise<string> {
    await new Promise((resolve) => setTimeout(() => resolve(true), 300));
    const res = await fetch(`${this.baseUrl}/bucket/upload-url?file=${name}`, {
      headers: { Authorization: `Bearer ${this.auth.getAccessToken()}` },
    });
    const data = (await res.json()) as SinteseResDTO<{ url: string }>;
    return data.data.url;
  }

  async uploadFile(file: File, url: string): Promise<void> {
    await new Promise((resolve) => setTimeout(() => resolve(true), 300));
    const res = await fetch(url, { method: "PUT", body: file });
    if (!res.ok) throw new Error("unable to upload file");
  }

  async createPhoto(photo: Partial<Photo>): Promise<ResponseDTO<Photo>> {
    const res = await fetch(`${this.baseUrl}/midia`, {
      headers: {
        Authorization: `Bearer ${this.auth.getAccessToken()}`,
        "Content-Type": "application/json",
      },
      body: JSON.stringify(photo),
      method: "POST",
    });
    const data = (await res.json()) as SinteseResDTO<SintesePhoto>;
    const p: Photo = {
      url: data.data.url,
      id: data.data.photoId,
      tags: ["tag"],
      title: data.data.title,
      createdAt: data.data.createdAt,
    };

    return {
      status: "success",
      data: p,
    };
  }
}
