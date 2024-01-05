// ignore_for_file: library_private_types_in_public_api

import 'package:flutter/material.dart';
import 'package:gouel/models/event_model.dart';
import 'package:gouel/services/gouel_api_service.dart';
import 'package:gouel/services/gouel_session_service.dart';
import 'package:gouel/widgets/gouel_button.dart';
import 'package:gouel/widgets/gouel_scaffold.dart';
import 'package:provider/provider.dart';

class EventsScreen extends StatefulWidget {
  const EventsScreen({super.key});

  @override
  _EventsScreenState createState() => _EventsScreenState();
}

class _EventsScreenState extends State<EventsScreen> {
  List<Event> events = [];

  @override
  void initState() {
    super.initState();
    _loadEvents();
  }

  Future<void> _loadEvents() async {
    try {
      var apiService = Provider.of<GouelApiService>(context, listen: false);
      var response = await apiService.getEvents(context);
      setState(() {
        events = response;
      });
    } catch (e) {
      // Gérer les erreurs ici
    }
  }

  @override
  Widget build(BuildContext context) {
    return GouelScaffold(
      body: Column(
        children: <Widget>[
          const Padding(
            padding: EdgeInsets.all(20),
            child: Text(
              'Choix de l\'événement',
              style: TextStyle(fontSize: 24, fontWeight: FontWeight.bold),
              textAlign: TextAlign.center,
            ),
          ),
          Expanded(
            child: ListView.separated(
              separatorBuilder: (context, index) {
                if (index < events.length - 1) {
                  return const SizedBox(
                    height: 12,
                  );
                }

                return const SizedBox();
              },
              itemCount: events.length,
              itemBuilder: (context, index) {
                var event = events[index];
                return GouelButton(
                    text: event.title,
                    onTap: () {
                      GouelSession().store("event", event);
                      Navigator.of(context).pushNamed(
                        '/manage_event',
                      );
                    });
              },
            ),
          ),
        ],
      ),
    );
  }
}
