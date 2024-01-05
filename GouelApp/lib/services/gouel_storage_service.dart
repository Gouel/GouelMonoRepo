import 'package:hive/hive.dart';

class GouelStorage {
  Box _openBox() {
    return Hive.box("gouelStorageBox");
  }

  Future<void> store(String key, dynamic value) async {
    var box = _openBox();
    await box.put(key, value);
  }

  Future<dynamic> retrieve(String key) async {
    var box = _openBox();
    return box.get(key);
  }

  Future<void> remove(String key) async {
    var box = _openBox();
    await box.delete(key);
  }
}
