import { useAppDispatch, useAppSelector } from "@store/store.ts";

import "./Upload.css";
import {
  getActionPhotoModalSet,
  getActionPhotosCreate,
  getActionUploadGetUrl,
  getActionUploadRequest,
} from "@store/photos/creators/photo.ts";
import { ActionStatus, ModalState } from "@core/domain/domain.ts";
import { ChangeEvent, FormEvent, useEffect, useState } from "react";
import * as url from "url";

const Modal = (opts: { children: React.ReactNode }) => {
  const photos = useAppSelector((state) => state.photos["photos/add"]);
  const dispatch = useAppDispatch();

  return (
    <div
      id={Math.random().toString(32).slice(2)}
      className="modal"
      style={{ display: photos.modal === ModalState.Open ? "block" : "none" }}
    >
      <div className="modal-content">
        <div className="modal-header">
          <span
            className="close"
            onClick={() => {
              dispatch(getActionPhotoModalSet(ModalState.Closed));
            }}
          >
            &times;
          </span>
          <h2>Novo</h2>
        </div>
        <div className="modal-body">{opts.children}</div>
      </div>
    </div>
  );
};

export const Upload = () => {
  const dispatch = useAppDispatch();
  const upload = useAppSelector((state) => state.photos["photos/upload"]);
  const create = useAppSelector((state) => state.photos["photos/add"]);
  const [file, setFile] = useState<File>();
  const [title, setTitle] = useState("");

  const handleChangeFile = (e: ChangeEvent<HTMLInputElement>) => {
    if (!e.target.files) return;
    setFile(e.target.files[0]);
  };

  const handleFormSubmit = (e: FormEvent) => {
    e.preventDefault();
    dispatch(
      getActionPhotosCreate({
        photo: {
          title: title,
          url: upload.upload_url?.split("?")?.[0] || "",
          fileName: file?.name,
        },
      })
    ).then(() => {
      setFile(undefined);
      setTitle("");
    });
  };

  useEffect(() => {
    if (!file) return;
    dispatch(getActionUploadGetUrl({ file }));
  }, [file]);

  useEffect(() => {
    if (!file || !upload.upload_url) return;
    dispatch(getActionUploadRequest({ file, url: upload.upload_url }));
  }, [upload.upload_url]);

  return (
    <Modal>
      <h2>Cadastro</h2>
      <form id="upload" onSubmit={handleFormSubmit}>
        <div className="input-wrapper">
          <label htmlFor="title">Título</label>
          <input
            type="title"
            id="title"
            placeholder="Tìtulo"
            value={title}
            onChange={(e) => setTitle(e.target.value)}
          />
        </div>
        <div className="input-wrapper-file">
          <label htmlFor="file">
            {(upload.status === ActionStatus.LOADING && "Carregando") ||
              (file && file.name) ||
              "Selecione o arquivo..."}
          </label>
          <input id="file" type="file" onChange={handleChangeFile} />
        </div>

        {upload.status === ActionStatus.ERROR && (
          <div className="error">
            <p>Erro ao realizar o carregamento da imagem</p>
          </div>
        )}
      </form>
      <div className="actions">
        <button
          form="upload"
          type="submit"
          disabled={!file || !url || upload.status === ActionStatus.LOADING}
        >
          {(create.status === ActionStatus.LOADING && "Loading...") || "Salvar"}
        </button>
      </div>
    </Modal>
  );
};
