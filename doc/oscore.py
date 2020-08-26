OntCversion = '2.0.0'
from ontology.interop.System.Storage import Put, GetContext, Get
from ontology.interop.System.Runtime import Notify,CheckWitness
from ontology.interop.Ontology.Contract import Migrate
from ontology.interop.System.Contract import Destroy
from ontology.libont import join

KeyPrefix = "SgSolo"
KeyOwnerAddress = "OwnerAddress"

def Main(operation, args):
    if operation == "init":
        return init(args)
    elif operation == "init_status":
        return init_status()
    elif operation == 'StoreUsedNum':
        return StoreUsedNum(args)
    elif operation == "GetUsedNum":
        return GetUsedNum(args)
    elif operation == "destroyContract":
        return DestroyContract()
    elif operation == "migrateContract":
        return MigrateContract(args)
    return False


def DestroyContract():
    addr = Get(GetContext(), KeyOwnerAddress)
    assert(len(addr) != 0)
    assert(CheckWitness(addr))
    return Destroy()


def MigrateContract(code):
    addr = Get(GetContext(), KeyOwnerAddress)
    assert(len(addr) != 0)
    assert(CheckWitness(addr))

    success = Migrate(code, True, "name", "version", "author", "email", "description")
    assert(success)
    Notify(["Migrate successfully", success])
    return success


def init_status():
    addr = Get(GetContext(), KeyOwnerAddress)
    Notify([addr])
    return addr


def init(addr):
    if len(Get(GetContext(), KeyOwnerAddress)) == 0:
        Put(GetContext(), KeyOwnerAddress, addr)
        Notify(["init True"])
        return True
    else:
        Notify(["init False"])
        return False


def StoreUsedNum(args):
    if len(args) != 3:
        return False

    addr = Get(GetContext(), KeyOwnerAddress)
    assert(len(addr) != 0)
    assert(CheckWitness(addr))

    userId = args[0]
    orderId = args[1]
    usedNum = args[2]
    t = concat(KeyPrefix,userId)
    KEY = concat(t, orderId)

    Put(GetContext(), KEY, usedNum)
    Notify([userId, orderId, usedNum])
    return True


def GetUsedNum(args):
    if len(args) != 2:
        return False

    userId = args[0]
    orderId = args[1]
    t = concat(KeyPrefix,userId)
    KEY = concat(t, orderId)

    num = Get(GetContext(), KEY)
    num = num + 0
    Notify([userId, orderId, num])
    return num
