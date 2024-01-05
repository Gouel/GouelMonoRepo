import random
import os


def generate_fake_events(num_events=10):
    titles = [
        "Concert de Rock",
        "Exposition d'Art Moderne",
        "Festival de Cinéma",
        "Conférence sur la Technologie",
        "Marathon de la Ville",
        "Salon Gastronomique",
    ]
    descriptions = [
        "Venez découvrir les talents locaux et internationaux.",
        "Une exposition captivante présentant les œuvres d'artistes contemporains.",
        "Projection des films les plus attendus de l'année.",
        "Des experts partagent leurs insights sur les dernières innovations.",
        "Participez ou encouragez les coureurs dans cette compétition emblématique.",
        "Dégustez des spécialités culinaires de renom.",
        "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Curabitur suscipit est nibh, ut congue erat porta non. Praesent posuere maximus tellus, eu tempus tortor posuere tincidunt. Suspendisse condimentum dolor vel ultrices laoreet. Suspendisse molestie, mi sed fringilla elementum, arcu risus tincidunt nisi, vitae tincidunt libero dui ac libero. Nullam purus quam, elementum vitae nulla quis, cursus tempus libero. Fusce sit amet gravida tortor, interdum viverra mauris. Morbi eu nunc odio. Mauris ac ipsum arcu. Quisque laoreet euismod nibh. Nunc varius elit sapien, non sollicitudin dolor luctus a. Phasellus rhoncus leo nec velit sollicitudin condimentum in ut lacus. Aliquam erat volutpat. Vivamus lobortis nulla at est eleifend vulputate. Interdum et malesuada fames ac ante ipsum primis in faucibus. ",
    ]

    locations = [
        "Salle Omega",
        "Galerie d'Art Luna",
        "Cinéma Central",
        "Centre de Conférence Orion",
        "Parc des Champions",
    ]
    contacts = [
        "info@omega.com",
        "contact@galerieluna.com",
        "films@cinemacentral.com",
        "tech@orionconf.com",
        "marathon@parcdeschampions.com",
    ]

    events = []
    for i in range(num_events):
        event = {
            "id": i + 1,
            "title": random.choice(titles),
            "description": random.choice(descriptions),
            "location": random.choice(locations),
            "contact": random.choice(contacts),
            "image_url": f"event{random.randint(1, 6)}.jpg",
        }
        events.append(event)

    return events


def generate_fake_tickets(num_tickets=5):
    ticket_names = [
        "Ticket Standard",
        "Ticket VIP",
        "Ticket Famille",
        "Ticket Étudiant",
        "Ticket Groupe",
    ]
    prices = [20, 50, 35, 15, 60]  # Exemple de prix

    tickets = []
    for i in range(num_tickets):
        ticket = {
            "id": i,
            "name": random.choice(ticket_names),
            "price": random.choice(prices),
        }
        tickets.append(ticket)

    return tickets


db = {
    "events": [
        {
            "public": True,
            "title": "Concert de Rock",
            "description": "Venez découvrir les talents locaux et internationaux.",
            "location": "Salle Omega",
            "contact": "info@omega.com",
            "image_url": f"event1.jpg",
            "tickets": [
                {
                    "name": "Place + Eco Cup",
                    "price": 3.0,
                },
                {
                    "name": "Place + Eco Cup + Collier collector",
                    "price": 8.0,
                },
            ],
        },
        {
            "public": True,
            "title": "Exposition d'Art Moderne",
            "description": "Une exposition captivante présentant les œuvres d'artistes contemporains.",
            "location": "Galerie d'Art Luna",
            "contact": "contact@galerieluna.com",
            "image_url": f"event2.jpg",
            "tickets": [
                {
                    "name": "Entrée",
                    "price": 5.0,
                },
                {
                    "name": "Entrée VIP",
                    "price": 25.0,
                },
            ],
        },
        {
            "public": True,
            "title": "Festival de Cinéma",
            "description": "Projection des films les plus attendus de l'année.",
            "location": "Cinéma Central",
            "contact": "films@cinemacentral.com",
            "image_url": f"event3.jpg",
            "tickets": [
                {
                    "name": "Place + Eco Cup",
                    "price": 3.0,
                },
                {
                    "name": "Place + Eco Cup + Collier collector",
                    "price": 8.0,
                },
            ],
        },
        {
            "public": True,
            "title": "Conférence sur la Technologie",
            "description": "Des experts partagent leurs insights sur les dernières innovations.",
            "location": "Centre de Conférence Orion",
            "contact": "tech@orionconf.com",
            "image_url": f"event4.jpg",
            "tickets": [
                {
                    "name": "Place + Eco Cup",
                    "price": 3.0,
                },
                {
                    "name": "Place + Eco Cup + Collier collector",
                    "price": 8.0,
                },
            ],
        },
        {
            "public": True,
            "title": "Marathon de la Ville",
            "description": "Participez ou encouragez les coureurs dans cette compétition emblématique.",
            "location": "Salle Omega",
            "contact": "info@omega.com",
            "image_url": f"event5.jpg",
            "tickets": [
                {
                    "name": "Place + Eco Cup",
                    "price": 3.0,
                },
                {
                    "name": "Place + Eco Cup + Collier collector",
                    "price": 8.0,
                },
            ],
        },
        {
            "public": True,
            "title": "Salon Gastronomique",
            "description": "Dégustez des spécialités culinaires de renom.",
            "location": "Parc des Champions",
            "contact": "marathon@parcdeschampions.com",
            "image_url": f"event6.jpg",
            "tickets": [
                {
                    "name": "Place + Eco Cup",
                    "price": 3.0,
                },
                {
                    "name": "Place + Eco Cup + Collier collector",
                    "price": 8.0,
                },
            ],
        },
    ],
    "users": [
        {
            "nom": "HARTMANN",
            "prenom": "Matthias",
            "email": "test@iziram.fr",
            "sessions": {},
            "solde": 5,
            "transactions": [
                {
                    "title": "Bière",
                    "price": 4.0,
                    "number": 1,
                    "date": "2023-12-07T22:30",
                },
                {
                    "title": "Soft",
                    "price": 2.0,
                    "number": 3,
                    "date": "2023-12-07T22:40",
                },
            ],
        },
        {
            "nom": "ONYMOUS",
            "prenom": "ANNE",
            "email": "test@iziram.fr",
            "sessions": {1: "benevole"},
            "solde": 5,
            "transactions": [
                {
                    "title": "Bière",
                    "price": 4.0,
                    "number": 1,
                    "date": "2023-12-07T22:30",
                },
                {
                    "title": "Soft",
                    "price": 2.0,
                    "number": 3,
                    "date": "2023-12-07T22:40",
                },
            ],
        },
        {
            "nom": "ADMIN",
            "prenom": "Admin",
            "email": "test@iziram.fr",
            "solde": 5,
            "sessions": {1: "admin"},
            "transactions": [
                {
                    "title": "Bière",
                    "price": 4.0,
                    "number": 1,
                    "date": "2023-12-07T22:30",
                },
                {
                    "title": "Soft",
                    "price": 2.0,
                    "number": 3,
                    "date": "2023-12-07T22:40",
                },
            ],
        },
        {
            "nom": "SUPERADMIN",
            "prenom": "SUPERADMIN",
            "email": "test@iziram.fr",
            "solde": 5,
            "sessions": {"superadmin": "admin"},
            "transactions": [
                {
                    "title": "Bière",
                    "price": 4.0,
                    "number": 1,
                    "date": "2023-12-07T22:30",
                },
                {
                    "title": "Soft",
                    "price": 2.0,
                    "number": 3,
                    "date": "2023-12-07T22:40",
                },
            ],
        },
    ],
    "benevoles": {
        0: [
            {"nom": "Bernard", "prenom": "Pierre", "droits": ["vestiaire", "buvette"]},
            {"nom": "Dubois", "prenom": "Pierre", "droits": ["vestiaire", "entrée"]},
            {
                "nom": "Dubois",
                "prenom": "Marie",
                "droits": ["entrée", "restauration", "prevention", "buvette"],
            },
            {
                "nom": "Robert",
                "prenom": "Sophie",
                "droits": ["buvette", "vestiaire", "caisse", "entrée"],
            },
            {
                "nom": "Bernard",
                "prenom": "Sophie",
                "droits": ["entrée", "restauration", "buvette", "caisse"],
            },
            {
                "nom": "Bernard",
                "prenom": "Luc",
                "droits": ["buvette", "entrée", "vestiaire", "prevention"],
            },
            {
                "nom": "Martin",
                "prenom": "Pierre",
                "droits": ["caisse", "restauration", "vestiaire", "buvette"],
            },
        ],
        1: [
            {"nom": "Thomas", "prenom": "Pierre", "droits": ["restauration", "entrée"]},
            {
                "nom": "Bernard",
                "prenom": "Marie",
                "droits": ["vestiaire", "entrée", "caisse"],
            },
            {"nom": "Martin", "prenom": "Sophie", "droits": ["entrée", "vestiaire"]},
            {"nom": "Thomas", "prenom": "Pierre", "droits": ["restauration"]},
            {
                "nom": "Thomas",
                "prenom": "Sophie",
                "droits": ["restauration", "vestiaire", "buvette", "entrée"],
            },
            {"nom": "Dubois", "prenom": "Marie", "droits": ["caisse"]},
            {
                "nom": "Robert",
                "prenom": "Sophie",
                "droits": ["entrée", "vestiaire", "buvette"],
            },
        ],
        2: [
            {
                "nom": "Robert",
                "prenom": "Luc",
                "droits": ["buvette", "caisse", "vestiaire"],
            },
            {
                "nom": "Bernard",
                "prenom": "Sophie",
                "droits": ["buvette", "caisse", "entrée", "vestiaire"],
            },
            {
                "nom": "Martin",
                "prenom": "Luc",
                "droits": ["vestiaire", "buvette", "restauration", "entrée"],
            },
            {
                "nom": "Thomas",
                "prenom": "Luc",
                "droits": ["buvette", "vestiaire", "caisse", "entrée"],
            },
            {"nom": "Dubois", "prenom": "Luc", "droits": ["restauration", "entrée"]},
        ],
        3: [
            {"nom": "Thomas", "prenom": "Sophie", "droits": ["buvette", "entrée"]},
            {
                "nom": "Martin",
                "prenom": "Sophie",
                "droits": ["entrée", "restauration", "caisse"],
            },
            {
                "nom": "Martin",
                "prenom": "Luc",
                "droits": ["restauration", "entrée", "buvette"],
            },
            {
                "nom": "Thomas",
                "prenom": "Luc",
                "droits": ["buvette", "caisse", "entrée", "vestiaire"],
            },
            {"nom": "Bernard", "prenom": "Luc", "droits": ["restauration", "buvette"]},
            {
                "nom": "Bernard",
                "prenom": "Sophie",
                "droits": ["buvette", "restauration"],
            },
            {
                "nom": "Thomas",
                "prenom": "Luc",
                "droits": ["vestiaire", "buvette", "entrée"],
            },
            {
                "nom": "Thomas",
                "prenom": "Jean",
                "droits": ["restauration", "vestiaire", "caisse", "entrée"],
            },
            {"nom": "Dubois", "prenom": "Pierre", "droits": ["buvette"]},
            {
                "nom": "Dubois",
                "prenom": "Luc",
                "droits": ["vestiaire", "caisse", "restauration", "entrée"],
            },
        ],
        4: [
            {
                "nom": "Robert",
                "prenom": "Marie",
                "droits": ["buvette", "restauration", "vestiaire", "caisse"],
            },
            {"nom": "Thomas", "prenom": "Jean", "droits": ["entrée", "vestiaire"]},
            {
                "nom": "Martin",
                "prenom": "Marie",
                "droits": ["restauration", "buvette", "vestiaire"],
            },
            {
                "nom": "Dubois",
                "prenom": "Pierre",
                "droits": ["restauration", "vestiaire"],
            },
            {
                "nom": "Martin",
                "prenom": "Jean",
                "droits": ["vestiaire", "buvette", "restauration", "entrée"],
            },
            {"nom": "Thomas", "prenom": "Luc", "droits": ["buvette", "caisse"]},
            {
                "nom": "Martin",
                "prenom": "Luc",
                "droits": ["caisse", "entrée", "vestiaire"],
            },
        ],
        5: [
            {
                "nom": "Bernard",
                "prenom": "Jean",
                "droits": ["entrée", "vestiaire", "caisse", "buvette"],
            },
            {
                "nom": "Robert",
                "prenom": "Marie",
                "droits": ["entrée", "caisse", "vestiaire", "buvette"],
            },
            {"nom": "Martin", "prenom": "Pierre", "droits": ["buvette", "caisse"]},
            {
                "nom": "Martin",
                "prenom": "Pierre",
                "droits": ["buvette", "vestiaire", "restauration", "caisse"],
            },
            {
                "nom": "Dubois",
                "prenom": "Pierre",
                "droits": ["entrée", "caisse", "buvette"],
            },
            {"nom": "Dubois", "prenom": "Luc", "droits": ["buvette", "vestiaire"]},
            {
                "nom": "Thomas",
                "prenom": "Jean",
                "droits": ["entrée", "buvette", "vestiaire"],
            },
            {
                "nom": "Robert",
                "prenom": "Sophie",
                "droits": ["buvette", "entrée", "vestiaire"],
            },
            {"nom": "Dubois", "prenom": "Sophie", "droits": ["vestiaire"]},
        ],
    },
    "tickets": [
        {
            "nom": "Dubois",
            "prenom": "Jean",
            "email": "admin@example.com",
            "valide": True,
        },
        {
            "nom": "Bernard",
            "prenom": "Pierre",
            "email": "info@example.com",
            "valide": False,
        },
        {"nom": "Dubois", "prenom": "Luc", "email": "info@example.com", "valide": True},
        {
            "nom": "Martin",
            "prenom": "Marie",
            "email": "contact@example.com",
            "valide": False,
        },
        {
            "nom": "Bernard",
            "prenom": "Sophie",
            "email": "contact@example.com",
            "valide": False,
        },
        {
            "nom": "Robert",
            "prenom": "Pierre",
            "email": "test@example.com",
            "valide": True,
        },
        {
            "nom": "Bernard",
            "prenom": "Pierre",
            "email": "admin@example.com",
            "valide": True,
        },
        {
            "nom": "Bernard",
            "prenom": "Jean",
            "email": "contact@example.com",
            "valide": True,
        },
        {
            "nom": "Martin",
            "prenom": "Marie",
            "email": "contact@example.com",
            "valide": True,
        },
        {
            "nom": "Thomas",
            "prenom": "Luc",
            "email": "contact@example.com",
            "valide": False,
        },
        {
            "nom": "Bernard",
            "prenom": "Luc",
            "email": "contact@example.com",
            "valide": True,
        },
        {
            "nom": "Thomas",
            "prenom": "Sophie",
            "email": "info@example.com",
            "valide": True,
        },
        {
            "nom": "Dubois",
            "prenom": "Marie",
            "email": "admin@example.com",
            "valide": False,
        },
        {
            "nom": "Robert",
            "prenom": "Marie",
            "email": "admin@example.com",
            "valide": True,
        },
        {
            "nom": "Dubois",
            "prenom": "Pierre",
            "email": "user@example.com",
            "valide": True,
        },
        {
            "nom": "Dubois",
            "prenom": "Jean",
            "email": "test@example.com",
            "valide": False,
        },
        {
            "nom": "Bernard",
            "prenom": "Marie",
            "email": "user@example.com",
            "valide": False,
        },
        {
            "nom": "Thomas",
            "prenom": "Luc",
            "email": "contact@example.com",
            "valide": True,
        },
        {
            "nom": "Robert",
            "prenom": "Marie",
            "email": "info@example.com",
            "valide": True,
        },
        {
            "nom": "Thomas",
            "prenom": "Jean",
            "email": "contact@example.com",
            "valide": True,
        },
    ],
}

if __name__ == "__main__":
    events = generate_fake_events()
    for event in events:
        print(event)
