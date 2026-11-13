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

**Service**:

- Register a doctor with valid details

**Database Queries**

- Insert a doctor record into a table
- Check for license number if it already exists

**Testing**

- Tested upon successful registration, duplicate license number, invalid doctor details

**Implemented Url(local)**

- (POST) http://localhost:8082/api/v1/doctors

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
