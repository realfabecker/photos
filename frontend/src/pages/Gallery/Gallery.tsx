import { useEffect, useState } from "react";
import { State, useAppDispatch, useAppSelector } from "@store/store.ts";
import { getActionLoadPhotoList } from "@store/photos/creators/photo.ts";
import Photo from "@pages/Gallery/Photo.tsx";
import { Photo as TPhoto } from "@core/domain/domain";
import { ActionStatus } from "@core/domain/domain.ts";

const ErrorSection = (opts: { message: string }) => {
  return (
    <>
      <h1>Galeria</h1>
      <section className="gallery">
        <p>{opts.message}</p>
      </section>
    </>
  );
};

const ActionsSection = (opts: {
  photos: State<TPhoto[]>;
  onClickLoadMore: () => void;
}) => {
  return (
    <div className="actions">
      <button
        disabled={opts.photos.status == ActionStatus.LOADING}
        onClick={opts.onClickLoadMore}
      >
        {opts.photos.status == ActionStatus.LOADING
          ? "Carregando..."
          : "Carregar mais fotos"}
      </button>
    </div>
  );
};

export default function Gallery() {
  const dispatch = useAppDispatch();
  const [params, setParams] = useState({ page: 1, limit: 6 });
  const photos = useAppSelector((state) => state.photos["photos/list"]);

  useEffect(() => {
    dispatch(
      getActionLoadPhotoList({
        page: params.page,
        limit: params.limit,
      })
    );
  }, [dispatch, params]);

  function handleClickLoadMore() {
    setParams((prevState) => ({
      ...prevState,
      page: prevState.page + 1,
    }));
  }

  if (photos.status === ActionStatus.ERROR) {
    return <ErrorSection message={photos?.error?.message || "Ops!"} />;
  }

  let i = 0;
  return (
    <>
      <h1>Galeria</h1>

      <section className="gallery">
        {(photos.data || []).map((p) => {
          const delay = 0.4 + i * 0.2;
          i = i >= params.limit ? 0 : i + 1;
          return <Photo key={p.id} photo={p} delay={delay} />;
        })}
      </section>

      <ActionsSection photos={photos} onClickLoadMore={handleClickLoadMore} />
    </>
  );
}
