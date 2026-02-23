import { mkdir, writeFile } from "node:fs/promises";
import { join } from "node:path";

function getType(db: string, ctype: string): string {
    if (ctype.startsWith("text(")) {
        const size = ctype.split("(")[1].split(")")[0];
        if (db === "mssql") return `NVARCHAR(${size})`;
        if (db === "mysql") return `VARCHAR(${size})`;
        return db === "sqlite" || db === "pg" ? "TEXT" : ctype;
    }
    
    const mapping: Record<string, Record<string, string>> = {
        "uuid": { sqlite: "TEXT", pg: "UUID", mysql: "CHAR(36)", mssql: "UNIQUEIDENTIFIER" },
        "text": { sqlite: "TEXT", pg: "TEXT", mysql: "TEXT", mssql: "NVARCHAR(MAX)" },
        "boolean": { sqlite: "INTEGER", pg: "BOOLEAN", mysql: "TINYINT(1)", mssql: "BIT" },
        "datetime": { sqlite: "INTEGER", pg: "TIMESTAMP", mysql: "DATETIME", mssql: "DATETIME2" },
        "binary(1024)": { sqlite: "BLOB", pg: "BYTEA", mysql: "VARBINARY(1024)", mssql: "VARBINARY(1024)" },
        "i32": { sqlite: "INTEGER", pg: "INTEGER", mysql: "INT", mssql: "INT" },
        "int": { sqlite: "INTEGER", pg: "INTEGER", mysql: "INT", mssql: "INT" },
        "json": { sqlite: "TEXT", pg: "JSONB", mysql: "JSON", mssql: "NVARCHAR(MAX)" }
    };
    
    return mapping[ctype]?.[db] || ctype;
}

type Column = [string, string, string];

function mkTable(db: string, name: string, cols: Column[], pk: string | undefined, foreignKeys: Record<string, any>): string {
    const lines = [`CREATE TABLE ${name} (`];
    const colDefs: string[] = [];
    
    for (const [cname, ctype, attrs] of cols) {
        const t = getType(db, ctype);
        let line = `    ${cname} ${t}`;
        
        if (cname === pk) {
            if (db === "sqlite" && t === "INTEGER") {
                line += " PRIMARY KEY AUTOINCREMENT";
            } else if (db === "mysql" && t === "INT") {
                line += " PRIMARY KEY AUTO_INCREMENT";
            } else if (db === "mssql" && t === "INT") {
                line += " PRIMARY KEY IDENTITY(1,1)";
            } else if (db === "pg" && (t === "INTEGER" || t === "SERIAL")) {
                line = `    ${cname} SERIAL PRIMARY KEY`;
            } else {
                line += " PRIMARY KEY";
            }
        } else if (attrs.includes("pk")) {
            // composite, handled later
        } else {
            if (!attrs.includes("nil")) {
                line += " NOT NULL";
            }
        }
        colDefs.push(line);
    }
    
    if (pk === undefined) {
        const pks = cols.filter(c => c[2].includes("pk")).map(c => c[0]);
        if (pks.length > 0) {
            colDefs.push(`    PRIMARY KEY (${pks.join(", ")})`);
        }
    }
    
    if (foreignKeys[name]) {
        const fks = foreignKeys[name];
        if (Array.isArray(fks) && Array.isArray(fks[0])) {
            for (const [col, refTable, refCol] of fks) {
                colDefs.push(`    FOREIGN KEY (${col}) REFERENCES ${refTable}(${refCol}) ON DELETE CASCADE`);
            }
        } else {
            const [col, refTable, refCol] = fks as [string, string, string];
            colDefs.push(`    FOREIGN KEY (${col}) REFERENCES ${refTable}(${refCol}) ON DELETE CASCADE`);
        }
    }

    lines.push(colDefs.join(",\n"));
    lines.push(");");
    return lines.join("\n");
}

function mkIndex(db: string, table: string, cols: string[], unique = false, name?: string, where?: string): string {
    name = name || `ix_${table}_${cols.join("_")}`;
    const uniqStr = unique ? " UNIQUE" : "";
    const whereStr = where ? ` WHERE ${where}` : "";
    return `CREATE${uniqStr} INDEX ${name} ON ${table} (${cols.join(", ")})${whereStr};`;
}

const schemas: Record<string, Column[]> = {
    "users": [
        ["id", "uuid", "pk"],
        ["primary_email", "text(256)", "nil"],
        ["primary_email_upcase", "text(256)", "nil"],
        ["primary_phone", "text(32)", "nil"],
        ["name", "text(256)", "nil"],
        ["name_upcase", "text(256)", "nil"],
        ["image_uri", "text(1024)", "nil"],
        ["is_banned", "boolean", ""],
        ["created_at", "datetime", ""],
        ["updated_at", "datetime", "nil"]
    ],
    "users_emails": [
        ["id", "uuid", "pk"],
        ["user_id", "uuid", "fk"],
        ["email", "text(256)", "nil"],
        ["email_upcase", "text(256)", "nil"],
        ["email_upcase_digest", "text(512)", ""],
        ["is_active", "boolean", ""],
        ["is_verified", "boolean", ""],
        ["is_primary", "boolean", ""],
        ["created_at", "datetime", ""],
        ["updated_at", "datetime", "nil"]
    ],
    "users_phones": [
        ["id", "uuid", "pk"],
        ["user_id", "uuid", "fk"],
        ["phone", "text(32)", "nil"],
        ["phone_digest", "text(32)", ""],
        ["is_active", "boolean", ""],
        ["is_verified", "boolean", ""],
        ["is_primary", "boolean", ""],
        ["created_at", "datetime", ""],
        ["updated_at", "datetime", "nil"]
    ],
    "user_password_auth": [
        ["user_id", "uuid", "pk,fk"],
        ["password_digest", "text(1024)", ""],
        ["password_expires_at", "datetime", "nil"],
        ["last_attempted_at", "datetime", "nil"],
        ["attempt_count", "int", ""],
        ["otp_digest", "text(1024)", "nil"],
        ["otp_expires_at", "datetime", "nil"],
        ["otp_link_token", "text(512)", "nil"],
        ["is_locked", "boolean", ""],
        ["created_at", "datetime", ""],
        ["updated_at", "datetime", "nil"]
    ],
    "user_sessions": [
        ["id", "uuid", "pk"],
        ["user_id", "uuid", ""],
        ["expires_at", "datetime", ""],
        ["token", "text", ""],
        ["ip_address", "text(40)", "nil"],
        ["user_agent", "text(256)", "nil"],
        ["created_at", "datetime", ""],
        ["updated_at", "datetime", "nil"]
    ],
    "user_totp": [
        ["id", "uuid", "pk"],
        ["user_id", "uuid", "fk"],
        ["secret_encrypted", "text(512)", ""],
        ["recovery_codes_hash", "text", ""],
        ["is_active", "boolean", ""],
        ["created_at", "datetime", ""],
        ["updated_at", "datetime", "nil"]
    ],
    "user_login_providers": [
        ["id", "uuid", "pk"],
        ["user_id", "uuid", "fk"],
        ["provider_id", "text(64)", ""],
        ["account_id", "text(128)", ""],
        ["access_token", "text", ""],
        ["refresh_token", "text", "nil"],
        ["id_token", "text", "nil"],
        ["access_token_expires_at", "datetime", "nil"],
        ["refresh_token_expires_at", "datetime", "nil"],
        ["scopes", "text", "nil"],
        ["created_at", "datetime", ""],
        ["updated_at", "datetime", "nil"]
    ],
    "user_passkey": [
        ["id", "uuid", "pk"],
        ["user_id", "uuid", "fk"],
        ["credential_id", "binary(1024)", ""],
        ["data", "json", ""]
    ],
    "user_claims": [
        ["id", "i32", "pk"],
        ["user_id", "uuid", "fk"],
        ["type", "text(256)", ""],
        ["value", "text", ""]
    ],
    "roles": [
        ["id", "uuid", "pk"],
        ["name", "text(64)", ""],
        ["name_upcase", "text(64)", ""],
        ["description", "text(512)", ""]
    ],
    "users_roles": [
        ["user_id", "uuid", "pk"],
        ["role_id", "uuid", "pk"]
    ],
    "role_claims": [
        ["id", "i32", "pk"],
        ["role_id", "uuid", "fk"],
        ["type", "text", ""],
        ["claim_value", "text", ""]
    ],
    "user_api_keys": [
        ["id", "i32", "pk"],
        ["user_id", "uuid", "fk"],
        ["name", "text(64)", ""],
        ["name_upcase", "text(64)", ""],
        ["key_hint", "text(16)", ""],
        ["key_digest", "text(1024)", ""],
        ["expires_at", "datetime", "nil"],
        ["is_locked", "boolean", ""],
        ["is_revoked", "boolean", ""],
        ["comment", "text(512)", "nil"],
        ["created_at", "datetime", ""],
        ["updated_at", "datetime", "nil"]
    ],
    "user_api_keys_roles": [
        ["user_api_key_id", "int", "pk"],
        ["role_id", "uuid", "pk"]
    ]
};

const foreignKeys: Record<string, any> = {
    "users_emails": ["user_id", "users", "id"],
    "users_phones": ["user_id", "users", "id"],
    "user_password_auth": ["user_id", "users", "id"],
    "user_sessions": ["user_id", "users", "id"],
    "user_totp": ["user_id", "users", "id"],
    "user_login_providers": ["user_id", "users", "id"],
    "user_passkey": ["user_id", "users", "id"],
    "user_claims": ["user_id", "users", "id"],
    "users_roles": [["user_id", "users", "id"], ["role_id", "roles", "id"]],
    "role_claims": ["role_id", "roles", "id"],
    "user_api_keys": ["user_id", "users", "id"],
    "user_api_keys_roles": [["user_api_key_id", "user_api_keys", "id"], ["role_id", "roles", "id"]]
};

async function main() {
    const dbs = ["sqlite", "pg", "mysql", "mssql"];
    
    for (const db of dbs) {
        const up: string[] = [];
        const down: string[] = [];
        
        for (const [table, cols] of Object.entries(schemas)) {
            const pk = cols.find(c => c[2].includes("pk") && !c[2].includes("fk") && !["users_roles", "user_api_keys_roles", "user_password_auth"].includes(table))?.[0];
            up.push(mkTable(db, table, cols, pk, foreignKeys));
            down.unshift(`DROP TABLE ${table};`);
            
            if (cols.some(c => c[0] === "user_id") && table !== "users_roles") {
                up.push(mkIndex(db, table, ["user_id"]));
            }
            if (cols.some(c => c[0] === "role_id") && table !== "users_roles" && table !== "user_api_keys_roles") {
                up.push(mkIndex(db, table, ["role_id"]));
            }
        }

        up.push(mkIndex(db, "users", ["primary_email_upcase"], true));
        up.push(mkIndex(db, "users", ["name_upcase"], true));
        up.push(mkIndex(db, "roles", ["name_upcase"], true));
        up.push(mkIndex(db, "user_api_keys", ["user_id", "name_upcase"], true));
        
        const isTrue = db === "sqlite" || db === "mssql" ? "1" : "true";
        if (db !== "mysql") {
            up.push(mkIndex(db, "users_emails", ["user_id"], true, "uq_users_emails_primary", `is_primary = ${isTrue}`));
            up.push(mkIndex(db, "users_phones", ["user_id"], true, "uq_users_phones_primary", `is_primary = ${isTrue}`));
        }

        const outDir = db === "pg" ? "postgres" : db;
        const baseDir = join("apps/api/db/migrations", outDir);
        await mkdir(baseDir, { recursive: true });
        
        await writeFile(join(baseDir, "000001_identity_schema.up.sql"), up.join("\n\n") + "\n");
        await writeFile(join(baseDir, "000001_identity_schema.down.sql"), down.join("\n") + "\n");
    }
    
    console.log("Migrations generated.");
}

main().catch(console.error);
