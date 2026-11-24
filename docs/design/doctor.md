# Doctor

- path: internal/doctor/

## Implementation:

**Doctor:**

- ID(uuid), Name, Speciality, Address, License Number, Department of type string
- CreatedAt, UpdatedAt of type time

**RegisterDoctorRequest**: (_user input_)

- Name (required)(length 2 - 160)
- Specialty (required)(length 2-80)
- Address (required)(length 2-80)
- License number (required and unique)(length 2-40)
- Department (required)(length 2- 80)

**Services**:

- Register a doctor with valid details
- Get doctor by Id
- Search a doctor by name
- Update doctor

**Database Queries**

- Insert a doctor record into a table
- Check for license number if it already exists
- Get doctor by doctor id
- Search doctor by name
- Get total count of doctors present by name

**Testing**

- Tested upon successful registration, duplicate license number, invalid doctor details, get request

**Implemented Url(local)**

- (POST) http://localhost:8082/api/v1/doctors
- (GET) http://localhost:8082/api/v1/doctors/{id}
- (GET) http://localhost:8082/api/v1/doctors?name={name}&page={page}&limit={limit}
- (PUT) http://localhost:8082/api/v1/doctors/{id}

### API responses:

**Doctor Register Response**

```
{
    "id": "2127632a-0632-4f63-93f9-387940d1625f",
    "name": "Dr. Prasad",
    "speciality": "Cardiologist",
    "address": "Mysore",
    "license_number": "kmc345",
    "department": "Cardiology",
    "created_at": "2026-11-13T10:11:17.969016Z",
    "updated_at": "2026-11-13T10:11:17.969016Z"
}

```

**Duplicate License Number**

```
{
    "timestamp": "2026-11-13T04:54:18927589173Z",
    "status": 409,
    "error": "Conflict",
    "message": "Duplicate licensenumber not allowed",
    "path": "/api/v1/doctors"
}
```

**Without Login**

```
{
    "timestamp": "2026-11-13T10:08:30.503929621Z",
    "status": 401,
    "error": "Unauthorized",
    "message": "Authentication is required",
    "path": "/api/v1/doctors"
}

```

**Register Doctor As a Patient**

```
{
    "timestamp": "2026-11-13T10:47:30.717151811Z",
    "status": 403,
    "error": "Forbidden",
    "message": "You do not have permission to perform this action",
    "path": "/api/v1/doctors"
}
```

**Get Doctor By Id**

```
{
    "id": "97002fe1-cd18-40c5-8ca2-5c4771a70cfc",
    "name": "Dr. Meena",
    "speciality": "Gastroenterologist",
    "address": "Mysore",
    "license_number": "kmc134",
    "department": "Gastroenterology",
    "created_at": "2026-11-18T06:21:28.827837Z",
    "updated_at": "2026-11-18T06:21:28.827837Z"
}
```

**Get Doctor without Id**

```
{
    "timestamp": "2026-11-18T07:29:00.562280843Z",
    "status": 400,
    "error": "Bad Request",
    "message": "Doctor ID is required",
    "path": "/api/v1/doctors"
}
```

**Get Doctor with wrong Id"**

```
{
    "timestamp": "2026-11-18T07:53:27.363068092Z",
    "status": 404,
    "error": "Not Found",
    "message": "Doctor not found",
    "path": "/api/v1/doctors"
}
```

**Search a Doctor with name in page 1**

```
{
    "doctor_list": [
        {
            "id": "61fac9bb-a681-4f9a-a83f-bd730c68bfd2",
            "name": "Dr. Jameela",
            "speciality": "Orthopedician",
            "address": "Mysore",
            "license_number": "kmc336",
            "department": "Orthopediology",
            "created_at": "2026-11-19T12:34:39.915774Z",
            "updated_at": "2026-11-19T12:34:39.915774Z"
        },
        {
            "id": "97002fe1-cd18-40c5-8ca2-5c4771a70cfc",
            "name": "Dr. Meena",
            "speciality": "Gastroenterologist",
            "address": "Mysore",
            "license_number": "kmc134",
            "department": "Gastroenterology",
            "created_at": "2026-11-18T06:21:28.827837Z",
            "updated_at": "2026-11-18T06:21:28.827837Z"
        }
    ],
    "total": 3
}
```

**Search a Doctor with name in page 2**

```
{
    "doctor_list": [
        {
            "id": "52d1322e-78d8-4053-a668-684e1916db2a",
            "name": "Dr. Meera",
            "speciality": "General Physician",
            "address": "Mysore",
            "license_number": "kmc144",
            "department": "General Medicine",
            "created_at": "2026-11-19T12:33:30.160371Z",
            "updated_at": "2026-11-19T12:33:30.160371Z"
        }
    ],
    "total": 3
}
```

**Doctor Update Response**

```
{
    "id": "81d5432c-bcea-443f-b3b9-df2d6ed25b91",
    "name": "Dr. Karan",
    "speciality": "Dental specialist",
    "address": "Mysore",
    "license_number": "kmc123",
    "department": "Dental",
    "created_at": "2026-11-13T04:16:29.696948Z",
    "updated_at": "2026-11-24T04:54:19.366009Z"
}
```
