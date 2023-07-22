import json
import plyvel
import argparse


def decode(key, value):
    """decode data"""
    key = key.decode('utf-8', errors='ignore')
    value = value.decode('utf-8', errors='ignore')
    value = json.loads(value)
    if key in ("settings", "state"):
        value = json.loads(value)
    return key, value


def parser(db, output_path):
    data = {}
    for key, value in db.iterator():
        key, value = decode(key, value)
        data[key] = value
    with open(output_path, "w", encoding="utf-8") as f:
        json.dump(data, f, ensure_ascii=False, indent=4)


if __name__ == "__main__":
    arg_parser = argparse.ArgumentParser()
    arg_parser.add_argument("-i", required=True, help="Input LevelDB path")
    arg_parser.add_argument("-o", default="tabs.json",
                            help="Output JSON file path")
    args = arg_parser.parse_args()
    db_path, output_path = args.i, args.o
    # open levelDB
    db = plyvel.DB(db_path)
    # parser data
    parser(db, output_path)
