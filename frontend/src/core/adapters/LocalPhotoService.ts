import "reflect-metadata";
import { injectable } from "inversify";
import { IPhotoService } from "@core/ports/ports.ts";
import { PagedDTO, Photo, ResponseDTO } from "@core/domain/domain.ts";

@injectable()
export class LocalPhotoService implements IPhotoService {
  constructor(private readonly storage = sessionStorage) {}
  async fetchPhotos(opts: {
    page: number;
    limit: number;
    token?: string;
  }): Promise<ResponseDTO<PagedDTO<Photo>>> {
    const data: Photo[] = JSON.parse(this.storage.getItem("_local") || "[]");
    const p = opts.page == 1 ? 0 : (opts.page - 1) * opts.limit;
    const items = data.reverse().slice(p, p + opts.limit);
    return {
      status: "success",
      data: {
        items: items,
        page_count: items.length,
        has_more: opts.page * opts.limit < data.length,
      },
    };
  }

  async createPhoto(photo: Partial<Photo>): Promise<ResponseDTO<Photo>> {
    const items: Record<string, Photo> = JSON.parse(
      this.storage.getItem("_upload") || "{}"
    );
    const target = items[photo?.fileName || ""];
    if (!target) {
      return Promise.reject("unable to obtain image file");
    }
    photo.url = target.url;
    photo.id = target.id;
    photo.fileName = target.fileName;
    photo.createdAt = new Date().toISOString();
    const data: Partial<Photo>[] = JSON.parse(
      this.storage.getItem("_local") || "[]"
    );
    data.push(photo);
    this.storage.setItem("_local", JSON.stringify(data));
    return {
      status: "success",
      data: photo as Photo,
    };
  }

  async getUploadUrl(name: string): Promise<string> {
    return Math.random().toString(32).slice(2) + name;
  }

  async uploadFile(file: File, url: string): Promise<void> {
    const dataURL = await new Promise((resolve, reject) => {
      const reader = new FileReader();
      reader.onload = () => resolve(reader.result);
      reader.onerror = reject;
      reader.readAsDataURL(file);
    });
    const items: Record<string, Partial<Photo>> = JSON.parse(
      this.storage.getItem("_upload") || "{}"
    );
    items[file.name] = {
      id: url,
      url: dataURL as string,
      fileName: file.name,
    };
    this.storage.setItem("_upload", JSON.stringify(items));
  }
}
