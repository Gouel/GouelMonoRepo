// ignore_for_file: use_build_context_synchronously

import 'package:flutter/material.dart';
import 'package:gouel/services/gouel_session_service.dart';
import 'package:gouel/services/gouel_storage_service.dart';
import 'package:gouel/utils/gouel_exception.dart';
import 'package:http/http.dart' as http;
import 'dart:convert';

class GouelRequest {
  final String _path;
  final String _method;
  final GouelStorage storage = GouelStorage();

  GouelRequest._(this._path, this._method);

  static GouelRequest get(String path) {
    return GouelRequest._(path, "GET");
  }

  static GouelRequest put(String path) {
    return GouelRequest._(path, "PUT");
  }

  static GouelRequest post(String path) {
    return GouelRequest._(path, "POST");
  }

  static GouelRequest delete(String path) {
    return GouelRequest._(path, "DELETE");
  }

  static Map<String, String> getHeaders(context) {
    GouelSession session = GouelSession();
    var token = session.retrieve("token");
    return {
      'Authorization': 'Bearer $token',
      'Content-Type': 'application/json',
    };
  }

  Future<dynamic> send(BuildContext context,
      {dynamic data, bool noHeaders = false}) async {
    var baseUrl = await storage.retrieve("server_addr");
    if (baseUrl == null || baseUrl == "") {
      throw GouelException(
          message: "Veuillez définir l'adresse du serveur",
          state: GouelExceptionState.critical);
    }

    var url = Uri.parse('$baseUrl$_path');
    http.Response response;
    final headers = noHeaders
        ? <String, String>{'Content-Type': 'application/json'}
        : getHeaders(context);

    switch (_method) {
      case 'PUT':
        response =
            await http.put(url, headers: headers, body: json.encode(data));
        break;
      case 'POST':
        response =
            await http.post(url, headers: headers, body: json.encode(data));
        break;
      case 'DELETE':
        response = await http.delete(url, headers: headers);
        break;
      case 'GET':
      default:
        response = await http.get(url, headers: headers);
    }
    final body = json.decode(response.body);
    if (response.statusCode == 200) {
      return body;
    } else {
      var m = body as Map<String, dynamic>;
      throw GouelException(
          message: "${m["error"]}", state: GouelExceptionState.critical);
    }
  }
}
