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

**Database Queries**

- Insert a doctor record into a table
- Check for license number if it already exists
- Get doctor by doctor id

**Testing**

- Tested upon successful registration, duplicate license number, invalid doctor details, get request

**Implemented Url(local)**

- (POST) http://localhost:8082/api/v1/doctors
- (GET) http://localhost:8082/api/v1/doctors/{id}

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
