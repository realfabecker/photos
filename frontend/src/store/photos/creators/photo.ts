import { createAction, createAsyncThunk } from "@reduxjs/toolkit";
import { Container } from "inversify";
import { IPhotoService, Types } from "@core/ports/ports.ts";
import { PhotosActions } from "@store/photos/actions/photos.ts";
import { Photo } from "@core/domain/domain.ts";

export const getActionLoadPhotoList = createAsyncThunk(
  "photos/list",
  async (opts: { page: number; limit: number; token?: string }, { extra }) => {
    const container = (<any>extra).container as Container;
    const service = container.get<IPhotoService>(Types.PhotoService);
    return service.fetchPhotos({
      page: opts.page,
      limit: opts.limit,
      token: opts.token,
    });
  }
);

export const getActionPhotoModalSet = createAction(
  "photos/modal_set",
  (state: string) => {
    return {
      type: PhotosActions.PHOTO_UPLOAD_VIEW,
      payload: { state },
    };
  }
);

export const getActionUploadGetUrl = createAsyncThunk(
  "photos_upload/get_url",
  async (opts: { file: File }, { extra }) => {
    const container = (<any>extra).container as Container;
    const service = container.get<IPhotoService>(Types.PhotoService);
    return service.getUploadUrl(opts.file.name);
  }
);

export const getActionUploadRequest = createAsyncThunk(
  "photos_upload/request",
  async (opts: { file: File; url: string }, { extra }) => {
    const container = (<any>extra).container as Container;
    const service = container.get<IPhotoService>(Types.PhotoService);
    return service.uploadFile(opts.file, opts.url);
  }
);

export const getActionPhotosCreate = createAsyncThunk(
  "photos/add",
  async (opts: { photo: Partial<Photo> }, { dispatch, extra }) => {
    const container = (<any>extra).container as Container;
    const service = container.get<IPhotoService>(Types.PhotoService);
    await service.createPhoto(opts.photo);
    dispatch(getActionLoadPhotoList({ page: 1, limit: 3 }));
  }
);
