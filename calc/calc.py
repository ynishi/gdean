import json
import sys

import pandas as pd


def max_emv(p1, ds1, ds2):
    df = pd.DataFrame([ds1, ds2])
    ps = pd.Series([p1, 1 - p1])
    emv = df.mul(ps, axis=0).sum()
    m = emv.idxmax()
    return {"ans": int(m)}


if __name__ == "__main__":
    ret = ""
    if len(sys.argv) == 1:
        print(ret)
        sys.exit()
    payload = json.loads(sys.argv[1])
    if payload.get("name") == "max_emv":
        data = payload.get("data")
        if data is not None:
            ret = max_emv(data.get("p1"), data.get("ds1"), data.get("ds2"))
    print(json.dumps(ret))
