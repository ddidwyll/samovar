package session

import (
  "samovar/lib/db"

  "ergo.services/ergo/gen"
)

type repo struct{ db.DbActor }

func newRepo() gen.ProcessBehavior {
  return &repo{}
}

func (r *repo) Init(_ ...any) error {
  r.InitDb("session", "[session_repo]")
  r.SetName("test")
  r.tryToRestore()
  return nil
}

func (r *repo) tryToRestore() {
  if data, err := r.Load(); err == nil {
    resp, err := r.Call(gen.Atom("session_state"), data)
    r.Log().Info("session.repo.tryToRestore.data: %#v, %s", resp, err)
  }
}
